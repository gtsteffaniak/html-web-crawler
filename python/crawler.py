import http.client
from urllib.parse import urlparse, urljoin
from html.parser import HTMLParser


class LinkParser(HTMLParser):
    def __init__(self, target_id=None):
        super().__init__()
        self.content = None
        self.target_id = target_id
        self.links = []
        self.in_target_element = False

    def handle_data(self, data):
        if self.target_id is None:
            if self.content is None:
                self.content = ""
            self.content += data

    def handle_starttag(self, tag, attrs):
        if tag == "a" and self.in_target_element:
            for attr in attrs:
                if attr[0] == 'href':
                    self.links.append(attr[1])
        elif self.target_id is not None:
            for attr in attrs:
                if attr[0] in ("id", "class") and attr[1] == self.target_id:
                    self.in_target_element = True
                    break  # Only need to check for the first matching attribute

    def handle_endtag(self, tag):
        if tag == self.target_id and self.target_id is not None:
            self.in_target_element = False


def fetch_html(url):
    print("fetching", url)
    parsed_url = urlparse(url)
    connection = http.client.HTTPSConnection(parsed_url.netloc)
    connection.request('GET', parsed_url.path or '/')
    response = connection.getresponse()
    if response.status == 200:
        content_type = response.getheader('Content-Type')
        if 'text/html' in content_type:
            return response.read().decode('utf-8')
        else:
            # If content is not HTML, return None
            return None
    else:
        return None


def crawl(url, target_id=None):
    domain = urlparse(url).netloc
    html_content = fetch_html(url)
    if html_content is None:
        return None

    parser = LinkParser(target_id)
    parser.feed(html_content)

    if target_id is None:
        return parser.content

    valid_links = []
    for link in parser.links:
        full_url = urljoin(url, link)
        parsed_full_url = urlparse(full_url)
        if parsed_full_url.netloc == domain and parsed_full_url.scheme == 'https':
            valid_links.append(full_url)
    pages_content = {}
    for link in list(set(valid_links)):
        if link in pages_content:
            continue
        pages_content[link] = fetch_html(link)

    return pages_content


# Example usage
url = "https://apnews.com/hub/earthquakes"
target_id = "PageList-items-item"  # Empty string to fetch the entire content
crawled_data = crawl(url, target_id)

if target_id is None:
    print("Fetched content:", crawled_data)
else:
    print("Total:", len(crawled_data))
