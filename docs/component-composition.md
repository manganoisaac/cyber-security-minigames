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

pages.EmailClient()
└─ layouts.Base
   └─ components.EmailList(data.Emails, selectedIndex, selectedCount)
      ├─ EmailItem(index, email, selected) ×N
      │  └─ EmailAvatar(index, email.Subject)
      ├─ EmailWord(word, selected, wordIndex, emailIndex) ×words
      └─ EmailSubmit(label, selectedCount, oob, emailIndex)
```

## Component Spec

| Name | File | Signature | Purpose |
|---|---|---|---|
| `Base` | `views/layouts/base.templ` | `Base(title string)` | Global layout: HTML document shell, stylesheet/script tags, and a `{ children... }` slot rendered into `<body hx-boost="true">`. |
| `Landing` | `views/pages/landing.templ` | `Landing()` | Landing page. Wraps the landing menu in `layouts.Base`. |
| `ChapterSelect` | `views/pages/chapter_select.templ` | `ChapterSelect()` | Chapter selection page. Wraps the chapter book and links in `layouts.Base`. |
| `EmailClient` | `views/pages/email_client.templ` | `EmailClient(emails []types.Email, selectedIndex int, wordsSelected int)` | Email minigame page. Wraps the email section in `layouts.Base` and passes `emails`, `selectedIndex`, and `wordsSelected` into `EmailList`. |
| `EmailList` | `views/components/email/email_list.templ` | `EmailList(emails []types.Email, selectedIndex int, selectedCount int)` | Two-pane email layout. Renders an `EmailItem` per email, an `EmailWord` per word of the selected email's content, and the `EmailSubmit` button. |
| `EmailItem` | `views/components/email/email_item.templ` | `EmailItem(index int, email types.Email, selected bool)` | A single email row in the sidebar. Renders an `EmailAvatar` plus truncated title/subject/content text. |
| `EmailAvatar` | `views/components/email/email_avatar.templ` | `EmailAvatar(index int, subject string)` | Circular avatar showing the subject's first letter with a color derived from the index. |
| `EmailWord` | `views/components/email/email_word.templ` | `EmailWord(word string, selected bool, i int, emailIndex int)` | A clickable word in the selected email's content; renders the `--selected` modifier variant when selected. `emailIndex` is threaded into the word-toggle URLs. |
| `EmailSubmit` | `views/components/email/email_submit.templ` | `EmailSubmit(label string, selectedCount int, oob bool, emailIndex int)` | Submit button showing the selected-word count. `oob=true` emits `hx-swap-oob` so word toggles update the badge out-of-band; `false` (the full `EmailList` render) keeps it inside the container swap when switching emails. `emailIndex` is threaded into the submit URL. |
