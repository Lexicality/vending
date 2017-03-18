package web

import (
	"context"
	"fmt"
	"net/http"

	"github.com/martini-contrib/render"
	logging "github.com/op/go-logging"
	stripe "github.com/stripe/stripe-go"
	"github.com/stripe/stripe-go/charge"

	"github.com/lexicality/vending/backend"
	"github.com/lexicality/vending/hardware"
	"github.com/lexicality/vending/vend"
)

// VendSession represents an attempt to vend, ideally so you can go find it again
type vendSession struct {
	ID    string
	State vend.Result
}

// StartVending starts a blocking attempt to vend
func (s *vendSession) startVending(
	ctx context.Context,
	hw hardware.Machine,
	stock backend.Stock,
	item *backend.StockItem,
) {
	stock.ReserveItem(ctx, item.ID)

	res := hw.Vend(ctx, item.Location)
	s.State = res

	stock.UpdateItem(context.TODO(), item.ID, res)
}

func handleBuy(
	req *http.Request,
	r render.Render,
	log *logging.Logger,

	stock backend.Stock,
	hw hardware.Machine,
) {
	reqCtx := req.Context()
	globalCtx := reqCtx.Value(globalContextKey).(context.Context)

	// TODO: Parse & validate
	itemID := req.PostFormValue("item")
	user := req.PostFormValue("stripeEmail")
	token := req.FormValue("stripeToken")

	item, err := stock.GetItem(reqCtx, itemID)
	if err != nil {
		log.Errorf("Unable to retrieve item %s: %s", itemID, err)
		r.HTML(500, "500", nil)
		return
	} else if item == nil {
		log.Warningf("Got payment request from %s for non-existant item %s", user, itemID)
		r.Text(400, "??? Missing item")
		return
	}

	params := &stripe.ChargeParams{
		Amount:   item.Price,
		Currency: "gbp",
		Desc:     item.Name,
	}
	params.AddMeta("user", user)
	params.SetSource(token)

	// At this point reqCtx is inaproprate since this needs to continue even if the user closes the page
	// However we probably don't want to stop *now* since money is happening.
	charge, err := charge.New(params)

	if err != nil {
		log.Debugf("OMG ERROR %+v %s", err, err)
		r.Text(500, err.Error())
		return
	}

	if globalCtx.Err() != nil {
		// um
		log.Criticalf("Bailing out of incomplete vend due to context closing: %s", globalCtx.Err())
		r.Text(503, "Please contact the trustees")
		return
	}

	// https://i.imgur.com/mibus.jpg
	vs := &vendSession{
		ID:    charge.ID,
		State: vend.NoResult,
	}
	go vs.startVending(
		globalCtx,
		hw,
		stock,
		item,
	)

	r.Text(200, fmt.Sprintf("%+v", charge))
}
