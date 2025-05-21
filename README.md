# go-Link 🔗

A lightweight CLI tool written in Go that checks web pages for broken links.

![image](https://github.com/user-attachments/assets/a20f56fd-b2ec-4161-b0e2-4e4cf36560a5)

## Features ✨

- Detects broken links (404, 500 errors, etc.)
- Handles both absolute and relative URLs
- Skips non-HTTP links (mailto:, tel:, etc.)
- Simple color-coded output (✅/❌)
- Fast link checking using HTTP HEAD requests
- No external dependencies (pure Go standard library)


## Examples
```bash
# Check a website
go run main.go https://tailwindcss.com/

# Check with protocol prefix automatically added
go run main.go tailwindcss.com
```
## To run 
```bash
go run main.go [URL]
```
## How It Works 🔧
- Fetches the HTML content of the given URL
- Extracts all `<a href="">` links
- Converts relative URLs to absolute URLs
- Checks each link's HTTP status code
- Reports broken links (non-200 status codes)

## Limitations ⚠️
- Doesn't handle JavaScript-rendered links
- May miss some dynamically generated URLs
- Rate limits not implemented (be cautious with large sites)
