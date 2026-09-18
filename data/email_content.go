package data

import "cyber-security-minigames/types"

var Emails = []types.Email{
	{
		ID:        "mobile-plan-renewal-notice",
		Title:     "Mobile Plan Renewal Notice from ISP",
		Subject:   "Mobile Plan Renewal Notice",
		To:        "Angela Bennett (angelabennett@metacortex.com.au)",
		From:      "Your ISP Provider (noreply@messages.yourispprovider.com.au)",
		Content:   "Hi Angela,\nThis is a reminder that we will process your monthly mobile phone plan payment for the number 0491 570 006 on 26/08/2027.\nThe amount charged will be $15 from your chosen payment method.\nFor more information, to cancel your plan, or to request more time to pay, please visit the Services -> Payments section of the Your ISP Provider app or website.\nWebsite:\n www.yourispprovider.com.au/myaccount/services/payments\nYou can also visit your local branch and speak to our staff.\nThis is an automated message. Please do not reply to this email.",
		Reason:    "It discusses money and contains a website address, which means it is still reasonable to examine it carefully. However, it does not pressure Angela to use a specific login link or provide sensitive information, and it gives her several independent ways to access the service, including manually opening the provider's app or website. The sender domain appears consistent with the organisation, although a correct-looking sender address should always be considered alongside the rest of the message rather than treated as proof by itself.",
		Timestamp: "1 Aug",
	},
	{
		ID:        "restaurant-order-confirmation",
		Title:     "Restaurant Order from Table Confirmation",
		Subject:   "Steak and Grill Order Confirmation - Order 29",
		To:        "Angela Bennett (angelabennett@metacortex.com.au)",
		From:      "Steak and Grill (noreplyorders@steakandgrill.com.au)",
		Content:   "Thank you for your order, Angela. We're preparing your meal now!\nOrder Number: 29\nOrdered: 20 August 2027 at 19:25\nTable: 18\nYour Order\n1 x Rump Steak - Medium Rare - $40\n1 x Large Thin Cut Chips - Add Gravy - $15\n1 x Lemon-Lime Bitters - $7\nTotal: $62\nSee a problem with your order? Please speak to one of our staff members at the counter.",
		Reason:    "The message contains specific information about an in-person transaction and does not ask Angela to click a link, download anything, provide information, or send money. It directs any concerns to staff physically present at the restaurant.",
		Timestamp: "2 Aug",
	},
	{
		ID:        "click-and-collect-order-ready",
		Title:     "Click and Collect Order Ready to Collect Notice",
		Subject:   "Your Order Is Ready to Collect - 93456",
		To:        "Angela Bennett (angelabennett@metacortex.com.au)",
		From:      "Cool Gadgets (info@coolgadgets.com.au)",
		Content:   "Hi Angela,\nYour order is ready to collect from Cool Gadgets Capalaba.\nPlease bring valid ID and your order number when collecting your items.\nYou have until 26 September 2027 to collect the order. If it is not collected by this date, the order will be cancelled and you will receive a full refund.\nOrder\nOrder Number: 93456\nDate Ordered: 15/08/2027\n1 x USB-C Charger 2m\n1 x 256GB microSD Card\n1 x Mini Action Camera\nTotal Paid: $425\nCollection Location\nCool Gadgets Capalaba\n16 Bay Street\nCapalaba QLD 4160\nIf you have questions about your order, please use the customer-support number listed on our official website or visit the store.\nPlease do not reply to this automated email.",
		Reason:    "There is a collection deadline, but the message does not threaten Angela with a financial loss. If she takes no action, she receives a refund. It does not ask for credentials or payment information and directs her to independently access the company's website or physical store if she needs assistance.",
		Timestamp: "3 Aug",
	},
	{
		ID:        "reminder-to-review-a-store",
		Title:     "Reminder to Review a Store",
		Subject:   "Let Us Know How We Did",
		To:        "Angela Bennett (angelabennett@metacortex.com.au)",
		From:      "Fun Fashion (noreply@funfashion.com.au)",
		Content:   "Dear valued customer,\nYour feedback is important to us.\nWe are always looking to improve your experience shopping with Fun Fashion, and you can help by leaving a short review.\nThe process should only take a minute.\nJust select how many stars you want to give us! Please note that these reviews are for internal use only and will not be shared publicly. We look forward to hearing your feedback.\nClick Here to Leave a Review\nDisplayed link:\n www.funfashion.com.au/reviews\nThank you for your time,\nFun Fashion",
		Reason:    "It uses a generic greeting, contains a link, and asks the recipient to interact with the email, but these do not automatically make an email malicious. The link and sender domain are consistent with the organisation, there is no request for sensitive information, and the message does not pressure Angela to act immediately. Many legitimate automated or marketing emails also use generic greetings.",
		Timestamp: "4 Aug",
	},
	{
		ID:                   "targeted-helios-access-review",
		Title:                "Helios Access Review",
		Subject:              "Project Helios Access Review",
		To:                   "Angela Bennett (angelabennett@metacortex.com.au)",
		From:                 "David Walker (david.walker@metac0rtex.com.au)",
		Content:              "Hi Angela,\nBefore you head off next week, we need to finish the access review for everyone currently assigned to Project Helios.\nIT has asked each team member to confirm their repository access so we don't run into any permission problems while you're away.\nPlease use the Helios Access Portal below and sign in with your MetaCortex account before 4:00pm today:\nReview Project Helios Access\nDisplayed address:\n https://portal.metacortex.com.au/helios\nActual link:\n https://metacortex-access-review.com/helios\nOnce you've confirmed your access, there's nothing else you need to do.\nThanks,\nDavid Walker\nProject Delivery Manager\nMetaCortex Industries",
		Reason:               "Eugene has used information gathered about Angela to make the message appear like a legitimate internal MetaCortex request. The sender domain uses a zero instead of the letter 'o' in metac0rtex, the displayed link points to the real MetaCortex portal but the actual destination is an unrelated external domain, and the email asks Angela to sign in with her MetaCortex account before a 4:00pm deadline. The accurate personal and workplace details make it believable, but they do not prove the sender is legitimate.",
		IsMalware:            true,
		Timestamp:            "5 Aug",
		CorrectPhraseIndexes: []int{2, 5, 7, 10, 12},
	},
	{
		ID:        "email-from-a-friend",
		Title:     "Email from a Friend",
		Subject:   "Can I Ask You a Favour?",
		To:        "Angela Bennett (angelabennett@metacortex.com.au)",
		From:      "Daniel White (danielwhite@myemail.com)",
		Content:   "Hi Angela,\nSorry to use your work email, but I know you use it for everything anyway and barely check your normal one.\nEllie's 40th is coming up next Sunday. I know you two don't really get along anymore and you weren't planning on going.\nI won't be able to go unless I can find someone to watch the dogs. Elvis is still a puppy and will destroy the house if nobody keeps an eye on him.\nI don't have much money, but I can pay you in that chocolate slice you always ask me to make.\nYou haven't been answering my texts either. Everything okay?\nGive me a call or text when you get a chance. You still have my number saved, right?\nDaniel",
		Reason:    "Using Angela's work email for a personal conversation is unusual, and the casual grammar is different from professional correspondence, but neither proves phishing. Most importantly, Daniel does not provide a link, an attachment, a request for credentials, or a new contact method. He asks Angela to contact him using the phone number she already has, giving her an independent way to verify the message.",
		Timestamp: "6 Aug",
	},
	{
		ID:        "sale-promo-code",
		Title:     "Sale Promo Code",
		Subject:   "Happy Birthday Angela - Treat Yourself On Us",
		To:        "Angela Bennett (angelabennett@metacortex.com.au)",
		From:      "Fun Fashion (promotions@funfashion.com.au)",
		Content:   "Hi Angela,\nHappy birthday!\nYou've received a $15 birthday coupon as part of your Fun Fashion Club membership.\nYou can access your coupon through the link below or by signing into your Fun Fashion account directly through our website.\nView Birthday Coupon\nDisplayed and actual address:\n www.funfashion.com.au/accounts/coupons/birthdayredeem\nYour coupon is valid for 12 months from the date of issue and can be used at participating Fun Fashion stores.\nThanks for being part of the Fun Fashion Club!",
		Reason:    "Rewards and free offers are commonly used in phishing, but this offer is relatively modest, there is no pressure to act quickly, the coupon remains valid for 12 months, and the sender and linked domains are consistent. Angela can obtain the coupon by independently signing into her account rather than using the email link. An incentive alone does not prove an email is malicious.",
		Timestamp: "7 Aug",
	},
	{
		ID:        "email-from-coworker",
		Title:     "Email from Coworker",
		Subject:   "IT Help Please",
		To:        "Angela Bennett (angelabennett@metacortex.com.au)",
		From:      "Robert Smith (robertsmith@metacortex.com.au)",
		Content:   "Morning Angela,\nHow are you going?\nI know I'm supposed to be going to Deb with the IT stuff but she's not getting back to me, so I wanted to ask you something.\nHow do I upload my document to the cloud again?\nEvery time I try it says something about not being logged in. I don't know, I'm hopeless with this stuff.\nI've attached a photo of the error I'm getting.\nCan you meet me outside my office tomorrow at lunch and show me what I'm doing wrong?\nThanks,\nRobert\nAttachment: cloud_login_error.jpg\nSent through MetaCortex internal communications.",
		Reason:    "It contains casual spelling and grammar, an unexpected attachment, and a request relating to account access, but none of these proves phishing. Poor grammar is common in informal communication, and an image attachment is lower risk than an unexpected executable, though any unexpected attachment should still be verified. The strongest reassuring factor is that Robert asks Angela to meet him in person at a known workplace location, so she can verify the request independently.",
		Timestamp: "8 Aug",
	},
	{
		ID:        "one-time-code-higher-risk",
		Title:     "One Time Code - Higher Risk",
		Subject:   "You Requested a One-Time Passcode",
		To:        "Angela Bennett (angelabennett@metacortex.com.au)",
		From:      "My Bank (donotreply@mybank.com.au)",
		Content:   "Dear Angela,\nA one-time passcode has been requested for your My Bank account.\nPasscode: 35FW23JYF\nThis request was made from:\nBrisbane, Australia\n13:00 AEST - 17/08/2027\nThis code will expire in five minutes.\nDo not share this code with anyone.\nIf you did not request this code, your account may be at risk.\nPlease open the official My Bank mobile app or visit our official website.\nFrom there, change your password and review your recent account activity.\nMy Bank will never ask you to send your password or one-time passcode through email.\nThis is an automated message. Please do not reply.",
		Reason:    "It contains urgency and a security warning, which are common in phishing, but genuine security alerts may also need users to act quickly. The important difference is how the email asks Angela to respond: it does not provide a login button or ask her to send information through email. Instead it instructs her to independently open an app she already has or manually enter the known website address. Urgency alone is not enough to classify an email as phishing.",
		Timestamp: "9 Aug",
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
