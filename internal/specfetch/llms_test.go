package specfetch

import (
	"slices"
	"testing"
)

const testIndex = `# Aikido Security Documentation

> Documentation for Aikido Security

## API Reference: Documentation
- [Authorization](https://docs.test/reference/authorization.md): Get started
  - [Verifying incoming webhooks](https://docs.test/reference/verify.md): Ensure

## API Reference: Authorization
- [Get access token](https://docs.test/reference/getaccesstoken.md): Retrieve a token

## API Reference: Issues
- [Export all issues](https://docs.test/reference/exportissues.md): Returns all issues
- [List open issue groups](https://docs.test/reference/listopenissuegroups.md): Use the header X-Has-Next-Page

## Something else
- [Blog](https://docs.test/blog/post.md): not an API page

## API Reference: Clouds
- [List connected clouds](https://docs.test/reference/listclouds-1.md): Returns clouds
`

func TestPageURLs(t *testing.T) {
	origin := mustURL(t, "https://docs.test/llms.txt")
	tests := []struct {
		name    string
		index   string
		want    []string
		wantErr string
	}{
		{
			name:  "keeps operation sections in index order",
			index: testIndex,
			want: []string{
				"https://docs.test/reference/exportissues.md",
				"https://docs.test/reference/listopenissuegroups.md",
				"https://docs.test/reference/listclouds-1.md",
			},
		},
		{
			name:  "keeps a title that contains brackets",
			index: "## API Reference: Issues\n- [[Beta] List things](https://docs.test/reference/beta.md): x\n",
			want:  []string{"https://docs.test/reference/beta.md"},
		},
		{
			name:    "rejects an entry it cannot read instead of skipping it",
			index:   "## API Reference: Issues\n- [Broken](https://docs.test/reference/broken): no .md\n",
			wantErr: `unrecognised index entry "- [Broken](https://docs.test/reference/broken): no .md"`,
		},
		{
			name:    "rejects a link to another host",
			index:   "## API Reference: Issues\n- [x](https://evil.test/reference/x.md): y\n",
			wantErr: "link https://evil.test/reference/x.md is not on https://docs.test",
		},
		{
			name:    "rejects plain http on an https origin",
			index:   "## API Reference: Issues\n- [x](http://docs.test/reference/x.md): y\n",
			wantErr: "is not on https://docs.test",
		},
		{
			name:    "fails when no operation pages are listed",
			index:   "## API Reference: Documentation\n- [a](https://docs.test/reference/a.md): b\n",
			wantErr: "index lists no API reference pages",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := PageURLs(tt.index, origin)
			if tt.wantErr != "" {
				wantErr(t, err, tt.wantErr)
				return
			}
			if err != nil {
				t.Fatal(err)
			}
			if !slices.Equal(got, tt.want) {
				t.Errorf("got  %q\nwant %q", got, tt.want)
			}
		})
	}
}
