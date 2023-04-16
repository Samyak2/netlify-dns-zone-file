# Netlify DNS Zone File Generator

Generate a zone file from your Netlify DNS records.

**Intended use case**: to transfer DNS records from Netlify to Cloudflare (which supports importing a zone file).

**How does it work**: Netlify does not provide an option to export records to a zone file, but it does provide an [API](https://open-api.netlify.com/#tag/dnsZone/operation/getDnsRecords) that lists them. This tool uses the data from this API and creates a zone file out of it.

## Usage

1. Create a [Netlify Personal Access Token](https://app.netlify.com/user/applications#personal-access-tokens).

Now, you can either [Run on Replit](https://replit.com/@samyaks/netlify-dns-zone-file) or run it locally on your system:

### Run on Replit

1. Export the netlify token as an environment variable:
    ```bash
    export NETLIFY_TOKEN=<your token here>
    ```
1. Run the tool. The output will contain the names of the `.zone` files that were generated.
    ```bash
    ./main
    ```

### Run locally

1. Clone this repository:
    ```bash
    git clone https://github.com/Samyak2/netlify-dns-zone-file.git
    ```
1. Export the netlify token as an environment variable:
    ```bash
    export NETLIFY_TOKEN=<your token here>
    ```
1. Run the tool. The output will contain the names of the `.zone` files that were generated.
    ```bash
    go run .
    ```


## License

MIT

## References

- [Wikipedia](https://en.wikipedia.org/wiki/Zone_file)
- [Netlify Open API docs](https://open-api.netlify.com/#tag/dnsZone/operation/getDnsRecords)
- [Importing a zone file in Cloudflare](https://developers.cloudflare.com/dns/manage-dns-records/how-to/import-and-export/#format-your-zone-file)
