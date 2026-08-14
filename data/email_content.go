package data

import "cyber-security-minigames/types"

var Emails = []types.Email{
	{
		ID:        "mobile-plan-renewal-notice",
		Title:     "Mobile Plan Renewal Notice from ISP",
		Subject:   "Mobile Plan Renewal Notice",
		To:        "Angela Bennett (angelabennett@metacortex.com.au)",
		From:      "Your ISP Provider (noreply@messages.yourispprovider.com.au)",
		Content:   "Hi Angela,\nThis is a reminder that we will process your monthly mobile phone plan payment for the number 0491 570 006 on the 26/08/27. The amount charged will be $15 from your chosen payment method. For more information, to cancel your plan, or to request more time to pay, please visit the Services -> Payments tab in the Your ISP Provider app or website (link leads to www.yourispprovider.com.au/myaccount/services/payments), or visit your local branch and speak to our friendly staff.\nThis is an automated message. Please do not reply to this email.",
		Reason:    "A legitimate reminder from your ISP with a correctly spelled sender address and safe, alternative contact methods.",
		Timestamp: "1 Aug",
	},
	{
		ID:        "restaurant-order-confirmation",
		Title:     "Restaurant Order from Table Confirmation",
		Subject:   "Steak and Grill Order Confirmation: 29",
		To:        "Angela Bennett (angelabennett@metacortex.com.au)",
		From:      "Steak and Grill (noreplyorders@steakandgrill.com.au)",
		Content:   "Thank you for your order, Angela. We’re prepping your grillicious meal now!\nHere is what we are cooking up for your order number 29:\n20th August 2027 Ordered at: 19:25\nTable Number: 18\n1 x Rump Steak - Medium Rare $40\n1 x Large Thin Cut Chips - Add Gravy $15\n1 x Lemon-Lime Bitters $7\nSubtotal: $62\nOrder Total: $62\nSee a problem with your order? Let one of our friendly staff know at the counter.",
		Reason:    "A legitimate receipt from a restaurant order with no urgency and no request for information.",
		Timestamp: "2 Aug",
	},
	{
		ID:        "click-and-collect-order-ready",
		Title:     "Click and Collect Order Ready to Collect Notice",
		Subject:   "Your Order is Ready to Collect - 93456",
		To:        "Angela Bennett (angelabennett@metacortex.com.au)",
		From:      "Cool Gadgets (info@coolgadgets.com.au)",
		Content:   "Hi Angela,\nYou have 3 items ready to collect from Cool Gadgets Capalaba. Please head to the front desk with your valid ID and order number ready to collect your items.\nYou have until the 26th of September 2027 to collect your items before we cancel your order and refund you the full amount.\nYour Order:\nOrder Number: 93456\nDate Ordered: 15/08/2027\n1 x USB C charger 2m\n1 x 256gb microSD card\n1 x Mini Action Camera\nSubtotal: $425\nTotal Paid: $425\nThis is not your receipt. A receipt was sent to this email on purchase and can also be collected from the counter on collection.\nWhere to Collect:\nCool Gadgets\nCapalaba\n16 Bay Street\n4160\nHours:\nMonday 8:30am - 5:30pm\nTuesday 8:30am - 5:30pm\nWednesday 8:30am - 5:30pm\nThursday 8:30am - 9:00pm\nFriday 8:30am - 5:30pm\nSaturday 9:30am - 5:00pm\nSunday 9:30am - 4:00pm\nIf you have any questions about your order, please do not hesitate to contact us using the 24 hour phone service provided on our website. We are always happy to help. Please do not reply to this email as this is an automated message.",
		Reason:    "A genuine collection notice; charges would be fully refunded if no action is taken, and contact is only via safe channels.",
		Timestamp: "3 Aug",
	},
	{
		ID:        "reminder-to-review-a-store",
		Title:     "Reminder to Review a Store",
		Subject:   "Let Us Know How We Did",
		To:        "Angela Bennett (angelabennett@metacortex.com.au)",
		From:      "Fun Fashion (noreply@funfashion.com.au)",
		Content:   "Dear valued customer,\nYour feedback is important to us. We are always looking to improve your experience shopping with us at Fun Fashion, and you can help us by leaving a review. The process is simple and will only take a minute of your time; all you have to do is select how many stars you want to give us! Please note that these reviews are for internal use only and will not be shared publicly. We look forward to hearing your feedback.\nClick here (link leads to www.funfashion.com.au/reviews) to access our review service.\nThank you for your time,\nFun Fashion",
		Reason:    "A harmless review request with no incentive, no urgency, and a correctly spelled sender address.",
		Timestamp: "3 Aug",
	},
	{
		ID:        "store-advertising-evil-hill-9",
		Title:     "Store Advertising a Newly Announced Product",
		Subject:   "Just Announced: Evil Hill 9",
		To:        "Angela Bennett (angelabennett@metacortex.com.au)",
		From:      "Games Corner (gamesnews@gamescorner.com.au)",
		Content:   "The evilest hill yet.\nJoin Jack Resident as he returns to that restless town, Evil Hill. This time, the towns are even more restless, and the hills are even more evil. Preorders open 15 September for SlayPlation 9 and YCube Zero. Get restless…get evil. “Evil Hill 9: What’s Going On With That Radio?”\nThis information is correct as of the date of the email campaign (5th September 2026). The Product is in development and may be subject to change. Product is subject to availability.",
		Reason:    "An unsolicited marketing email with no links to follow, no incentive, and no request for information.",
		Timestamp: "4 Aug",
	},
	{
		ID:        "email-from-coworker",
		Title:     "Email from Coworker",
		Subject:   "It help please help",
		To:        "Angela Bennett (angelabennett@metacortex.com.au)",
		From:      "Robert Smith (robertsmith@metacortex.com.au)",
		Content:   "Morning Angela, how are you going? I know that I;m supposed to be going to deb with the it stuff but shes not getting back to me so i just wanted to ask you a question. How do i upload my document to the cloud again? Everytime time i try i says something about not being loggin in I dont know man i just need a hand im not good with this tech stuff. Can you meet me outside my office tomorrow lunch time?\nEmail sent with jpeg attachment of a photo taken (presumably with a phone) of the failed login screen for the cloud service named something like rjrh4o439f398.\nSent from MetaCortex internal communications.",
		Reason:    "Sent from a colleague within MetaCortex; meeting in person is a safe contact and the jpeg attachment contains no executable code.",
		Timestamp: "5 Aug",
	},
	{
		ID:        "email-from-a-friend",
		Title:     "Email from a Friend",
		Subject:   "",
		To:        "Angela Bennett (angelabennett@metacortex.com.au)",
		From:      "Daniel White (danielwhite@myemail.com)",
		Content:   "Hi Angela,\nSorry to use your work email, but I know you use it for everything anyway and don’t check your normal one. Ellie’s 40th party is coming up soon. Next sunday. I know you and her don’t get on anymore, and I know you aren’t planning on going. I won’t be able to go unless I get someone to watch the dogs for me. Elvis is just a puppy and will destroy the whole house if he doesn’t have someone watching him 24/7. You know I don’t have much money, I can pay you with my famous chocolate slice you’ve always loved.\nYou haven’t been answering my texts… are you doing ok? Have I done something to upset you? I’d like to clear the air with you as well. Give me a call or a text or something when you can, ok? You still have me saved in your phone, right?\nKind regards,\nDaniel",
		Reason:    "An email from a friend asking for a favour in exchange for a minimal, non-monetary reward and no request for information.",
		Timestamp: "5 Aug",
	},
	{
		ID:        "store-sale-notification",
		Title:     "Store Sale Notification",
		Subject:   "SALE! SALE! SALE!",
		To:        "Angela Bennett (angelabennett@metacortex.com.au)",
		From:      "Camping Gear (news@campinggear.com.au)",
		Content:   "CAMPING GEAR’S ULTIMATE CLEARANCE SALE STARTS NEXT WEEK!\nDON’T MISS OUT ON THESE AWESOME DEALS! SAVE UP TO 80% RRP ON SELECTED ITEMS!\nONLY $149 FOR A WATERPROOF SINGLE-PERSON TENT! DOWN FROM $179!\nFISHING RODS STARTING AT $50! DOWN FROM $70!\nHALF PRICE ON SELECTED KAYAKS!\nIN STORE ONLY AND WHILE STOCKS LAST! DON’T MISS OUT!\nCAMPING GEAR - FOR ALL YOUR CAMPING NEEDS!",
		Reason:    "A poorly put together marketing campaign that requests no information and provides no way to buy through the email.",
		Timestamp: "5 Aug",
	},
	{
		ID:        "one-time-code-low-risk",
		Title:     "One Time Code - Low Risk",
		Subject:   "You requested a verification code",
		To:        "Angela Bennett (angelabennett@metacortex.com.au)",
		From:      "Books Galore (noreply@booksgalore.accounts.com.au)",
		Content:   "Hello Angela,\nHere is your verification code: 1358320284hdh3\nUse this code to verify that the email address tied to this account belongs to you.\nThis code expires in 15 minutes and can only be used once. Please don’t share this code. If you did not request this code, you can ignore this email. This is an automatic service, please do not reply.",
		Reason:    "A verification code email with no incentive, correct sender information, and no request for information.",
		Timestamp: "6 Aug",
	},
	{
		ID:        "one-time-code-higher-risk",
		Title:     "One Time Code - Higher Risk",
		Subject:   "One Time Password",
		To:        "Angela Bennett (angelabennett@metacortex.com.au)",
		From:      "My Bank (donotreply@mybank.com.au)",
		Content:   "Dear Angela,\nYou have requested a one time passcode: 35fw23jyf\nIf you did not request this passcode, your account may be at risk. Please contact us from the MyBank app or website, change your password, and check your account history for suspicious activity immediately. Always access MyBank by either using the MyBank app or typing exactly mybank.com.au into an internet browser search bar.\nThis login request was made from Brisbane, Australia at 13:00 17/08/2027 AEST.\nThis password is valid for 5 minutes. Do not share this code with anyone. Do not reply to this email.\nWarning: MyBank will never send you a link to log into your bank account or ask for financial/login information. For more information, please visit www.mybank.com.au/security",
		Reason:    "A bank code to be used via the safe app or exact URL, with a correctly spelled sender address and contact details.",
		Timestamp: "6 Aug",
	},
	{
		ID:        "sale-promo-code",
		Title:     "Sale Promo Code",
		Subject:   "Happy Birthday Angela - Treat Yourself On Us",
		To:        "Angela Bennett (angelabennett@metacortex.com.au)",
		From:      "Fun Fashion (promotions@funfashion.com.au)",
		Content:   "Hi Angela,\nCongratulations! You’ve earned a birthday coupon worth $15! To access your coupon click here (link leads to www.funfashion.com.au/accounts/coupons/birthdayredeem), or from your Fun Fashion account page on our website. Thanks for being a member of the Fun Fashion Club.\nCoupon is valid for 12 months after issue and is only valid at participating Fun Fashion stores.",
		Reason:    "A legitimate birthday coupon with a 12 month expiry (no urgency) and a correctly spelled sender address.",
		Timestamp: "7 Aug",
	},
	{
		ID:        "overdue-payment-paylater",
		Title:     "Overdue Payment with a Buy Now Pay Later Service",
		Subject:   "You Have an Overdue Payment",
		To:        "Angela Bennett (angelabennett@metacortex.com.au)",
		From:      "PayLater (noreply@paylater.com.au)",
		Content:   "Dear Angela,\nYou have 1 overdue payment of $243 AUD from Cool Gadgets Capalaba. While this payment is overdue, you will be unable to make further purchases on PayLater. All PayLater customers receive a 3 month grace period to pay overdue fees. We will charge you a late fee of 10% of the overdue payment for every month overdue starting on the 17/11/2027 if you do not make this payment on or before this time.\nThank you for your consideration.\nWarning: PayLater will never ask you to send us financial or personal information over email or send you links. Please access our services using the PayLater mobile app.\nFor more assistance, please contact our customer support line on our app. This is an automated email, please do not reply.",
		Reason:    "A genuine overdue notice directing contact through the safe PayLater app, with a correctly spelled sender address.",
		Timestamp: "7 Aug",
	},
	{
		ID:        "passkey-creation-notification",
		Title:     "Passkey Creation Notification Email",
		Subject:   "A new passkey has been created on your Company account",
		To:        "Angela Bennett (angelabennett@metacortex.com.au)",
		From:      "MetaCortex (donotreply@metacortex.com.au)",
		Content:   "Dear Angela Bennett,\nA new passkey has been created for your MetaCortex account. You can use this passkey to sign in securely across multiple devices without using a password.\nIf you did not create this passkey, please immediately sign into your MetaCortex account and change your password, remove any unrecognised passkeys, check account history for suspicious activity, and alert management.\nKind regards,\nMetaCotrex Support\nSent from MetaCortex internal communications.",
		Reason:    "Sent from within MetaCortex with correct sender information and instructions that only engage safe, in-app channels.",
		Timestamp: "7 Aug",
	},
	{
		ID:                   "fake-reward",
		Title:                "Fake Reward",
		Subject:              "Cool Gadgets - Claim Your Rewards",
		To:                   "Angela Bennett (angelabennett@metacortex.com.au)",
		From:                 "Cool Gadgets (info@coolcadgets.com.au)",
		Content:              "Dear valued customer,\nYou have 1204 unused rewards points, redeeminable for a value of $743.35 AUD.\nThis email is to notify you that these rewards points will expire in 24 hours.\nRedeem your rewards here: www.coolgadgets.com.au/rewards (actual link: www.coolcadgets.com.au/rewards).\nIf you have any questions, please do not hesitate to respond to this email.\nThank you for being part of our rewards system. We hope you enjoy your day.",
		Reason:               "Phishing: a too good to be true reward, a generic greeting, a misspelled sender address (coolcadgets), a hidden link to a misspelled web address, and urgency.",
		IsMalware:            true,
		Timestamp:            "8 Aug",
		CorrectPhraseIndexes: []int{2, 4, 5, 6, 7},
	},
}

func EmailByID(id string) (types.Email, int, bool) {
	for i, em := range Emails {
		if em.ID == id {
			return em, i, true
		}
	}
	return types.Email{}, 0, false
}

// PhishingEmail returns the email the player must find to progress to level 2.
func PhishingEmail() (types.Email, int, bool) {
	for i, em := range Emails {
		if em.IsMalware {
			return em, i, true
		}
	}
	return types.Email{}, 0, false
}
