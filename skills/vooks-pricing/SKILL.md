---
name: vooks-pricing
description: Answer what a Vooks subscription costs, including the verified-teacher rate, the free trial, gifting, and how schools and districts buy in bulk.
license: Vooks API terms, https://www.vooks.com/terms-of-service/
homepage: https://www.vooks.com/plans/
---

# What Vooks costs

Use this whenever someone asks the price of Vooks, compares plans, asks about a teacher
discount, or asks how a school buys it. Do not answer from memory: prices change, and
this endpoint is the published source.

## The endpoint

```
curl https://www.vooks.com/api/v1/plans
```

No API key. Returns every priced plan plus the offerings that are quoted rather than
priced.

```json
{
  "data": [
    { "id": "consumer-annual", "price": { "amount": 69.99, "currency": "USD" },
      "interval": "year", "audience": "consumer", "trial_days": 7 }
  ],
  "quoted": [
    { "id": "school-district", "pricing": "quote", "contact": "support@vooks.com" }
  ],
  "trial_days": 7
}
```

## The shape of the answer

- **Consumer**: annual or monthly. Annual works out materially cheaper per month.
- **Teacher**: the same library at a verified-teacher rate, annual or monthly. Teachers
  verify at <https://www.vooks.com/educators/>.
- **Trial**: every paid plan includes a 7-day free trial. Accounts activated by
  redeeming a code are not eligible for it.
- **Schools and districts**: no self-serve price. Quoted per customer through the form
  at <https://www.vooks.com/educators/> or by emailing support@vooks.com. Vooks is also
  available to schools through MackinVIA.
- **Gifting**: <https://www.vooks.com/gift-vooks/>. Gift codes are redeemed, so they are
  not eligible for the trial.

## Watch for

Subscriptions bought inside the iOS or Android app are billed by Apple or Google and are
managed in that store's subscription settings, not on vooks.com. Say so when someone
asks how to cancel.

The prose version, for quoting, is <https://www.vooks.com/pricing.md>.

## Related

- Searching the library: the `vooks-catalog` skill
- MCP tool: `get_pricing_plans` on `https://www.vooks.com/mcp`
