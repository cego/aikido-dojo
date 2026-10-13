# Automatic Top-Ups

Buy credits automatically when your wallet drops below a threshold, so credit-backed scans keep running.

Automatic top-ups let Aikido buy more credits when your wallet runs low. You set the balance that triggers a purchase, how many credits to add, and an optional monthly cap.

Automatic top-ups are not available by credit card. Aikido sends an invoice, and the credits are added after it is paid.

## What you need

* Workspace admin access.
* At least one successful payment on the workspace, by card or by invoice. A [manual credit purchase](/miscellaneous-info/wallet-and-credits.md#how-to-purchase-credits) counts. That payment unlocks automatic top-ups. Later top-ups are still invoiced.

## Turn it on

Open [Wallet settings](https://app.aikido.dev/settings/wallet/settings), click **Manage** next to **Automatic Top-Ups**, turn it on, and click **Save**.

| Setting         | Required | What it does                                                                                     |
| --------------- | -------- | ------------------------------------------------------------------------------------------------ |
| Top-Up Trigger  | Yes      | Aikido creates an invoice when the wallet balance drops below this many credits.                 |
| Top-Up Amount   | Yes      | Credits added each time. Minimum 100. One credit is $1 USD.                                      |
| Monthly Maximum | No       | The most credits Aikido invoices automatically in one calendar month. Leave it empty for no cap. |

The top-up amount cannot be higher than the monthly maximum.

## What happens when the wallet runs low

1. Aikido checks workspaces that have automatic top-ups turned on.
2. If the balance is below the trigger, Aikido creates an invoice for the top-up amount and emails it to workspace admins.
3. The purchase shows in [wallet history](https://app.aikido.dev/settings/wallet/history) as **Automatic Top-Up** with status **Awaiting Payment**.
4. After the invoice is paid, Aikido adds the credits and the transaction becomes **Completed**.

Aikido does not create another automatic top-up while an invoice is still unpaid. Pay or settle that invoice before the next one can be issued.

If the balance is still below the trigger after a top-up is paid, Aikido can invoice again, up to the monthly maximum. The monthly maximum resets at the start of each calendar month. It only counts automatic top-ups, not credits you buy yourself and not [monthly complimentary credits](/miscellaneous-info/wallet-and-credits.md#monthly-complimentary-credits).

## Spending limits

The monthly maximum caps how many credits Aikido buys for you. [Spending limits](/miscellaneous-info/wallet-and-credits.md#spending-limits) cap how many credits a feature can use. Set both if you want a ceiling on purchases and a ceiling on usage.
