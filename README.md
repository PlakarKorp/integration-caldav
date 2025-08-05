
# CalDAV Integration

## What is CalDAV?

**CalDAV** is a standard protocol for accessing and managing calendar data on a remote server. It allows clients to read, write, and sync calendar events using HTTP-based requests. CalDAV is supported by many calendar providers, including Nextcloud, Fastmail, Google, and Apple.

This integration allows:

- Seamless backup of calendar events from CalDAV servers into a Plakar repository
- Direct restoration of `.ics` calendar files to remote CalDAV destinations
- Compatibility with any CalDAV-compliant server or service


## Installation

If a pre-built package exists for your system and architecture,
you can simply install it using:

```sh
$ plakar pkg add caldav
```

Otherwise,
you can first build it:

```sh
$ plakar pkg build caldav
```

This should produce `caldav-vX.Y.Z.ptar` that can be installed with:

```bash
$ plakar pkg add ./caldav-v0.1.0.ptar
```

## Configuration

The configuration parameters are as follows:

- `location` (required): The URL of the CalDAV server in the form `caldav://<host>[:<port>]`
- `username` (required for most providers): The username to authenticate as
- `password` (required for most providers): The password or app-specific token to authenticate with

For providers that require OAuth2 (e.g. Google, Microsoft, Apple), you must also supply:

- `name`: The provider name (e.g. google, microsoft, apple)
- `client_id`: Your OAuth2 client ID
- `client_secret`: Your OAuth2 client secret
- `service_scope`: Provider-specific OAuth2 scope

These values must be set in the plugin configuration when adding the source or destination.

## Example Usage

```bash
# configure a CalDAV source
$ plakar source add mycaldav caldav://cal.example.org username=<user> password=<pass>

# backup the source (import events)
$ plakar backup @mycaldav

# configure a CalDAV destination
$ plakar destination add mycaldav caldav://cal.example.org username=<user> password=<pass>

# restore the snapshot to the destination (export events)
$ plakar restore -to @mycaldav <snapid>
```

## Backup Format

Calendar events are saved in standard **`.ics` (iCalendar)** format, including:

- Event title, start/end time, recurrence rules
- Attendees and organizer information
- Calendar and event metadata (UID, creation date, etc.)

## Round-Trip Support

> ✅ Events imported from CalDAV can be exported back without data loss. The format is preserved 1:1 across import/export operations to ensure full fidelity.

## Tips

- **Use app-specific passwords** for providers like iCloud or Fastmail that don’t allow regular account passwords.
- **Read-only mode?** If your account is restricted, export operations might fail — check permissions.
- **Filter support:** Currently, all accessible calendars are imported/exported. Per-calendar selection may be added later.

## Feedback

Spotted an issue or have a feature in mind?  
Open an issue on the Plakar repository or reach out to the team.

---

## 🚀 Usage

### ⬇️ Import calendar events into Plakar

```bash
plakar source add mycaldav caldav://<url> username=<user> password=<pass>
plakar at /path/to/repo backup @mycaldav
```

This fetches all calendar events accessible via the CalDAV endpoint and stores them as `.ics` files in the repository.

### ⬆️ Export calendar events from Plakar

```bash
plakar destination add mycaldav caldav://<url> username=<user> password=<pass>
plakar at /path/to/repo restore @mycaldav
```

This pushes `.ics` calendar files previously stored in Plakar back to your CalDAV server.

---

## 🔐 OAuth2 Configuration

Some CalDAV providers (e.g. Google, Microsoft, Apple, or enterprise platforms) require **OAuth2** instead of basic username/password authentication.

To use such services, you must manually supply the required OAuth2 fields in your Plakar configuration:

```yaml
name: <provider>             # e.g. google, microsoft, apple
client_id: <your-client-id>
client_secret: <your-client-secret>
service_scope: <provider-specific-scope>
```

These values must be set in the plugin configuration (`config`) when adding the source or destination.

Plakar uses them to retrieve an access token from the appropriate OAuth2 endpoint.

> The `name` field determines which OAuth2 provider configuration to use. Make sure it matches one of the supported providers built into the plugin.


---

## 📂 Backup Format

Calendar events are saved in standard **`.ics` (iCalendar)** format, including:

- Event title, start/end time, recurrence rules
- Attendees and organizer information
- Calendar and event metadata (UID, creation date, etc.)

---

## 🔄 Round-Trip Support

> ✅ Events imported from CalDAV can be exported back without data loss.

The format is preserved 1:1 across import/export operations to ensure full fidelity.

---

## 🛠️ Tips

- **Use app-specific passwords** for providers like iCloud or Fastmail that don’t allow regular account passwords.
- **Read-only mode?** If your account is restricted, export operations might fail — check permissions.
- **Filter support:** Currently, all accessible calendars are imported/exported. Per-calendar selection may be added later.

---

## 📬 Feedback

Spotted an issue or have a feature in mind?  
Open an issue on the Plakar repository or reach out to the team.
