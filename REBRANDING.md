# NekoBox → KiviBox Rebranding Guide

## Summary

This fork has been rebranded from "NekoBox" to "KiviBox" with a fully configurable application name system.

## Quick Start

To change the application name, simply edit `conf/app.ini`:

```ini
[app]
title = "YourAppName"
```

Then rebuild and restart the Docker container.

## Changes Made

### 1. Removed Chinese Changelog
- **File**: `templates/change-logs.html`
- **Change**: Disabled external changelog feed from original author
- The changelog page now shows "This page has been disabled"

### 2. Configurable App Name System

#### Configuration
- **File**: `conf/app.ini`
- **Setting**: `title = "KiviBox"`
- Change this value to customize the application name everywhere

#### Backend Changes
- **File**: `internal/conf/static.go`
  - Added `Title string \`ini:"title"\`` field to App struct
  
- **File**: `internal/context/context.go` (line 213-214)
  - Changed from hardcoded `"NekoBox"` to `conf.App.Title`
  - Added `c.Data["AppTitle"] = conf.App.Title` for templates

- **File**: `internal/mail/mail.go`
  - Updated all email subjects to use `conf.App.Title`
  - Updated email sender name: `"From", fmt.Sprintf("%s <%s>", conf.App.Title, conf.Mail.Account)`
  - Added `params["appTitle"] = conf.App.Title` for email templates

### 3. Template Updates

#### Web Templates
- `templates/base/header.html`: Logo uses `{{ .AppTitle }}`
- `templates/home.html`: Heading uses `{{ .AppTitle }}`
- `templates/base/footer.html`: Updated GitHub links to `EthanPany/NekoBox`

#### Email Templates
All email templates now use `{{ .appTitle }}` variable:
- `templates/mail/new-question.html`
- `templates/mail/new-answer.html`
- `templates/mail/password-recovery.html`

#### Other Templates
Replaced hardcoded "NekoBox" with "KiviBox":
- `templates/maintenance-mode.html`
- `templates/sponsor.html`
- `templates/user/profile.html`

### 4. Where the Name Appears

The configurable `title` value now controls:

1. **Web UI**:
   - Browser page title: `<title>{{ .AppTitle }}</title>`
   - Navigation bar logo: `<a class="uk-navbar-item uk-logo">{{ .AppTitle }}</a>`
   - Homepage heading: `<h1>{{ .AppTitle }}</h1>`

2. **Email**:
   - Email sender: `KiviBox <box@ethanpan.me>`
   - Email subjects: `【KiviBox】您有一个新的提问`
   - Email content: References throughout email body

3. **System**:
   - Default page title in context middleware
   - Flash messages and notifications

## Rebranding Your Fork

To rebrand this fork with your own name:

1. **Update the configuration**:
   ```bash
   # Edit conf/app.ini
   [app]
   title = "YourAppName"
   ```

2. **Update documentation** (optional):
   - Edit README.md if needed
   - Update DEPLOYMENT.md references
   - Update this REBRANDING.md

3. **Rebuild and restart**:
   ```bash
   docker build -t nekobox:latest .
   docker stop nekobox && docker rm nekobox
   docker run -d -p 8001:8080 --restart=always --name nekobox nekobox:latest
   ```

4. **Verify the changes**:
   ```bash
   curl http://localhost:8001 | grep "<title>"
   # Should show: <title>YourAppName</title>
   ```

## Benefits of This Approach

✅ **Single source of truth**: Change one config value to rebrand everywhere
✅ **No code changes needed**: Pure configuration-based rebranding
✅ **Fork-friendly**: Easy for others to customize their forks
✅ **Consistent branding**: Name appears consistently across all touchpoints
✅ **Maintainable**: Future updates won't break your branding

## Original Project

This is a fork of [NekoBox](https://github.com/wuhan005/NekoBox) by E99p1ant.

Original features and functionality remain intact - only branding has been modified.

## Technical Notes

- The `AppTitle` template variable is set in `internal/context/context.go:214`
- Email templates receive `appTitle` via params in `internal/mail/mail.go:56`
- The title is loaded from `conf.App.Title` which reads from `conf/app.ini`
- No hardcoded "NekoBox" strings remain in user-facing areas

## Testing Checklist

After rebranding, verify:
- [ ] Homepage shows correct title
- [ ] Browser tab shows correct title  
- [ ] Navigation bar logo is correct
- [ ] Email subjects use correct name
- [ ] Email sender name is correct
- [ ] Email body content is correct
- [ ] Maintenance mode page shows correct name

---

**Last Updated**: 2025-11-05
**Current Branding**: KiviBox
**Configuration**: `conf/app.ini`
