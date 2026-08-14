# Component Composition

The templ component tree for the Cyber Security Minigames application. Only templ components are shown; plain HTML elements are omitted.

## Composition Tree

```
layouts.Base(title)
└─ { children... }

pages.Landing()
└─ layouts.Base

pages.ChapterSelect()
└─ layouts.Base

pages.EmailClientPage(emails, active, selectedPhrases, level)
└─ layouts.Base
   └─ <section class="email-client"> + <div id="modal-root">
      └─ components.EmailClient(emails, active, selectedPhrases)
         ├─ EmailSidebar(emails, active)
         │  └─ EmailItem(email, active, oob) ×N
         │     └─ EmailAvatar(email)
         ├─ EmailDetails(email, selectedPhrases)
         │  └─ EmailPhrase(phrase, selected, phraseIndex, emailID) ×phrases
         └─ EmailSubmit(label, selectedCount, oob, emailID)
```

## Component Spec

| Name | File | Signature | Purpose |
|---|---|---|---|
| `Base` | `views/layouts/base.templ` | `Base(title string)` | Global layout: HTML document shell, stylesheet/script tags, and a `{ children... }` slot rendered into `<body hx-boost="true">`. |
| `Landing` | `views/pages/landing.templ` | `Landing()` | Landing page. Wraps the landing menu in `layouts.Base`. |
| `ChapterSelect` | `views/pages/chapter_select.templ` | `ChapterSelect()` | Chapter selection page. Wraps the chapter book and links in `layouts.Base`. |
| `EmailClientPage` | `views/pages/email_client.templ` | `EmailClientPage(emails []types.Email, active types.Email, selectedPhrases map[int]bool, level int)` | Email minigame page. Renders the `EmailClient` component inside `layouts.Base`, a level indicator, and the `#modal-root` target that verdict modals swap into. |
| `EmailClient` | `views/components/email/email_client.templ` | `EmailClient(emails []types.Email, active types.Email, selectedPhrases map[int]bool)` | Two-pane email layout. Renders an `EmailItem` per email via `EmailSidebar`, an `EmailPhrase` per phrase of the active email's content, and the `EmailSubmit` button. |
| `EmailSidebar` | `views/components/email/email_sidebar.templ` | `EmailSidebar(emails []types.Email, active types.Email)` | Renders one `EmailItem` per email, marking the active one with the `--active` modifier. |
| `EmailItem` | `views/components/email/email_item.templ` | `EmailItem(email types.Email, active bool, oob bool)` | A single email row in the sidebar. Renders an `EmailAvatar`, truncated title/subject/content text, and a checkbox that posts `/email-check/{id}` and swaps the verdict modal into `#modal-root`. `oob=true` emits `hx-swap-oob`. |
| `EmailAvatar` | `views/components/email/email_avatar.templ` | `EmailAvatar(email types.Email)` | Circular avatar showing the subject's first letter with a color derived from the email id via `avatar.Color`. |
| `EmailDetails` | `views/components/email/email_details.templ` | `EmailDetails(email types.Email, selectedPhrases map[int]bool)` | Renders the active email's content as clickable `EmailPhrase`s. |
| `EmailPhrase` | `views/components/email/email_phrase.templ` | `EmailPhrase(phrase string, selected bool, i int, emailID string)` | A clickable phrase in the selected email's content; renders the `--selected` modifier variant when selected. `emailID` is threaded into the phrase-toggle URLs. |
| `EmailSubmit` | `views/components/email/email_submit.templ` | `EmailSubmit(label string, selectedCount int, oob bool, emailID string)` | Submit button showing the selected-phrase count. `oob=true` emits `hx-swap-oob` so phrase toggles update the badge out-of-band; `false` (the full `EmailClient` render) keeps it inside the container swap when switching emails. `emailID` is threaded into the submit URL. |
| `Modal` | `views/components/ui/modal.templ` | `Modal(message string, url string)` | Overlay modal with a backdrop and an "Okay" button. Empty `url` dismisses via `DELETE /email-client/modal`; a non-empty `url` makes closing GET that URL with `hx-target="body"` so the page re-renders (used for the level-change refresh after finding the phishing email). |
