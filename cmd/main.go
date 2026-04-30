package main

import (
	"fmt"
	"maps"
	"net/url"
	"os"
	"os/exec"
	"strconv"
	"strings"

	"github.com/kissanjamgit/preview/brazz"
	"github.com/kissanjamgit/preview/config"
	"github.com/kissanjamgit/preview/hd8k"
	"github.com/kissanjamgit/preview/manyv"
	"github.com/kissanjamgit/preview/pbox"
	"github.com/kissanjamgit/preview/pworld"
	smex "github.com/kissanjamgit/preview/sexmex"
	"github.com/kissanjamgit/preview/teamskt"
	"github.com/kissanjamgit/preview/vip"

	"github.com/kissanjamgit/preview"

	"github.com/spf13/cobra"
)

///todo
//https://tours-store.psmcdn.net/swap_bundle/_search?sort=publishedDate:desc&q=(type:video%20AND%20isXSeries:false%20)&size=30&from=30
// GET /swap_bundle/_search?sort=publishedDate:desc&q=(type:video%20AND%20isXSeries:false%20)&size=30&from=30 HTTP/2
// Host: tours-store.psmcdn.net
// User-Agent: Mozilla/5.0 (Windows NT 10.0; Win64; x64; rv:150.0) Gecko/20100101 Firefox/150.0
// Accept: */*
// Accept-Language: en-US,en;q=0.9
// Accept-Encoding: gzip, deflate, br, zstd
// Referer: https://www.swappz.com/
// Origin: https://www.swappz.com
// Sec-GPC: 1
// Connection: keep-alive
// Sec-Fetch-Dest: empty
// Sec-Fetch-Mode: cors
// Sec-Fetch-Site: cross-site
// Priority: u=0
//https://tours-store.psmcdn.net/freeusebundle/_search?q=(site.seo.seoSlug.keyword:%22freeuse-fantasy%22%20AND%20type:video)&sort=publishedDate:desc&size=30&from=30
// GET /freeusebundle/_search?q=(site.seo.seoSlug.keyword:%22freeuse-fantasy%22%20AND%20type:video)&sort=publishedDate:desc&size=30&from=30 HTTP/2
// Host: tours-store.psmcdn.net
// User-Agent: Mozilla/5.0 (Windows NT 10.0; Win64; x64; rv:150.0) Gecko/20100101 Firefox/150.0
// Accept: */*
// Accept-Language: en-US,en;q=0.9
// Accept-Encoding: gzip, deflate, br, zstd
// Referer: https://www.freeuse.com/
// Origin: https://www.freeuse.com
// Sec-GPC: 1
// Connection: keep-alive
// Sec-Fetch-Dest: empty
// Sec-Fetch-Mode: cors
// Sec-Fetch-Site: cross-site
// Priority: u=0
// TE: trailers
//nubiles-porn
//https://en.inkasex.com/videos/latest
//https://pornbox.com/application/studio/328
//https://pornbox.com/application/model/188419
//https://www.moderndaysins.com/en/video/moderndaysins/One-Bed-Two-In-Laws/261172
///

func makeMapWithList[K comparable, V any](keys []K, defaultValue V) map[K]V {
	m := make(map[K]V, len(keys))
	for _, k := range keys {
		m[k] = defaultValue
	}
	return m
}

// func domainList() (l []string) {
// 	l = append(l, brazz.Domain...)
// 	l = append(l, pbox.Domain...)
// 	l = append(l, "teamskeet")
// 	l = append(l, "pornworld")
// 	l = append(l, "sexmex")
//
// 	return
// }

type Domain struct {
	domain string
	view   func(string) preview.Preview
}

func domainMap() map[string]func(string) preview.Preview {
	m := make(map[string]func(string) preview.Preview)
	brazzList := makeMapWithList(brazz.Domain, brazz.New)
	pList := makeMapWithList(pbox.Domain, pbox.New)
	maps.Copy(m, brazzList)
	maps.Copy(m, pList)
	m["teamskeet"] = teamskt.New
	m["pornworld"] = pworld.New
	m["sexmex"] = smex.New
	m["vip4k"] = vip.New
	m["pornhd8k"] = hd8k.New
	return m
}

func toDomainList(list []string, f func(string) preview.Preview) (d []Domain) {
	for _, i := range list {
		d = append(d, Domain{i, f})
	}
	return
}

func run() (pr []preview.ContentResource, err error) {
	// input := flag.Int("i", 0, "select the provider")
	// search := flag.String("search", "", "search for preview")
	// index := flag.Int("index", 0, "index")
	// help := flag.Bool("h", false, "print help text")
	// flag.Parse()

	// var domain []Domain
	// domain = append(domain, toDomainList(brazz.Domain, brazz.New)...)
	// domain = append(domain, toDomainList(pbox.Domain, pbox.New)...)
	// domain = append(domain, Domain{`teamskeet`, teamskt.New})
	// domain = append(domain, Domain{`pornworld`, pworld.New})
	// domain = append(domain, Domain{`sexmex`, smex.New})
	// domain = append(domain, Domain{`vip4k`, vip.New})
	// domain = append(domain, Domain{`pornhd8k`, hd8k.New})

	// if *search != "" {
	// 	var hostname string
	// 	if uri, err := url.Parse(*search); err != nil && uri.Hostname() != `` {
	// 		return nil, err
	// 	} else {
	// 		hostname = uri.Hostname()
	// 	}
	//
	// 	switch {
	// 	case strings.Contains(hostname, `manyvids`):
	// 		m := manyv.New(*search)
	// 		return m.Search(*index)
	// 	default:
	// 		return nil, fmt.Errorf("switch case exhausted")
	//
	// 	}
	// }

	// if *help || *input == -1 {
	// 	var buff strings.Builder
	// 	for i, d := range domain {
	// 		fmt.Fprintf(&buff, "%2d: %s\n", i, d.domain)
	// 	}
	// 	fmt.Println(buff.String())
	// 	return nil, nil
	// }

	// str := domain[*input].domain
	// pr, err = domainMap()[str](str).Get(*index)
	// if err != nil {
	// 	fmt.Println(err)
	// 	os.Exit(1)
	// }
	return
}

var domain = buildDomains()

func buildDomains() (d []Domain) {
	d = append(d, toDomainList(brazz.Domain, brazz.New)...)
	d = append(d, toDomainList(pbox.Domain, pbox.New)...)
	d = append(d, toDomainList(teamskt.Domain, brazz.New)...)

	d = append(d, Domain{`pornworld`, pworld.New},
		Domain{`sexmex`, smex.New},
		Domain{`vip4k`, vip.New},
		Domain{`pornhd8k`, hd8k.New},
	)
	return
}

var example = func() string {
	var buff strings.Builder
	for i, d := range domain {
		fmt.Fprintf(&buff, "%2d: %s\n", i, d.domain)
	}
	return buff.String()
}()

func cobraPlay(_ *cobra.Command, args []string, index int, cfg config.Config) (err error) {
	if len(args) == 0 {
		err = fmt.Errorf("need at least one argument")
		return
	}
	if len(args) > 1 {
		err = fmt.Errorf("too many arguments")
		return
	}

	input, err := strconv.Atoi(args[0])
	if err != nil {
		return
	}

	domain := domain[input]
	pr, err := domain.view(domain.domain).Get(index)
	if err != nil {
		return
	}

	var buffer strings.Builder
	buffer.WriteString("#EXTM3U\n")
	for _, cr := range pr {
		buffer.WriteString("#EXTINF:-1," + cr.Source + "\n" + cr.View + "\n")
	}
	PlayerArgs := []string{"-"} // stdin input
	PlayerArgs = append(PlayerArgs, strings.Split(cfg.PlayerArgs, " ")...)
	playCmd := exec.Command(cfg.Player, PlayerArgs...)

	stdin, err := playCmd.StdinPipe()
	if err != nil {
		return
	}
	err = playCmd.Start()
	if err != nil {
		return
	}
	_, err = stdin.Write([]byte(buffer.String()))
	if err != nil {
		return
	}
	err = stdin.Close()
	if err != nil {
		return
	}
	err = playCmd.Wait()

	return
}

func cobraShow(_ *cobra.Command, args []string, index int) (err error) {
	if len(args) == 0 {
		err = fmt.Errorf("need at least one argument")
		return
	}
	if len(args) > 1 {
		err = fmt.Errorf("too many arguments")
		return
	}

	input, err := strconv.Atoi(args[0])
	if err != nil {
		return
	}

	domain := domain[input]

	// pr, err := domainMap()[str](str).Get(index)
	pr, err := domain.view(domain.domain).Get(index)
	if err != nil {
		return
	}

	var buffer strings.Builder
	buffer.WriteString("#EXTM3U\n")
	for _, cr := range pr {
		buffer.WriteString("#EXTINF:-1," + cr.Source + "\n" + cr.View + "\n")
	}
	fmt.Println(buffer.String())
	return
}

func cli() (err error) {
	cfg, err := config.New()
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	cbr := cobra.Command{
		Use: "view",
	}
	var index int
	show := cobra.Command{
		Use:     "show",
		Example: example,
		RunE: func(cmd *cobra.Command, args []string) error {
			return cobraShow(cmd, args, index)
		},
	}

	play := cobra.Command{
		Use:     "play",
		Example: example,
		RunE: func(cmd *cobra.Command, args []string) error {
			return cobraPlay(cmd, args, index, cfg)
		},
	}
	search := cobra.Command{
		Use:     "search",
		Example: example,
		RunE: func(cmd *cobra.Command, args []string) (err error) {
			var hostname string
			if uri, err := url.Parse(args[0]); err != nil && uri.Hostname() != `` {
				return err
			} else {
				hostname = uri.Hostname()
			}

			switch {
			case strings.Contains(hostname, `manyvids`):
				m := manyv.New(args[0])
				m.Search(index)
			default:
				err = fmt.Errorf("switch case exhausted")
				return

			}
			return
		},
	}

	show.Flags().IntVar(&index, "index", 0, "index value")
	play.Flags().IntVar(&index, "index", 0, "index value")
	cbr.AddCommand(&show)
	cbr.AddCommand(&play)
	cbr.AddCommand(&search)
	err = cbr.Execute()
	if err != nil {
		os.Exit(1)
	}
	return
}

func main() {
	cli()
}
