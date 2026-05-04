package main

import (
	"fmt"
	"net/url"
	"os"
	"os/exec"
	"strconv"
	"strings"

	"github.com/kissanjamgit/preview/allanal"
	"github.com/kissanjamgit/preview/brazz"
	"github.com/kissanjamgit/preview/config"
	"github.com/kissanjamgit/preview/czechav"
	"github.com/kissanjamgit/preview/dorcel"
	"github.com/kissanjamgit/preview/enjoyx"
	"github.com/kissanjamgit/preview/evil"
	"github.com/kissanjamgit/preview/ftherapy"
	"github.com/kissanjamgit/preview/hd8k"
	"github.com/kissanjamgit/preview/lifese"
	"github.com/kissanjamgit/preview/littlecaprice"
	"github.com/kissanjamgit/preview/manyv"
	"github.com/kissanjamgit/preview/newsens"
	"github.com/kissanjamgit/preview/nporn"
	"github.com/kissanjamgit/preview/pbox"
	"github.com/kissanjamgit/preview/private"
	"github.com/kissanjamgit/preview/pworld"
	smex "github.com/kissanjamgit/preview/sexmex"
	"github.com/kissanjamgit/preview/teamskt"
	"github.com/kissanjamgit/preview/vip"

	"github.com/kissanjamgit/preview"

	"github.com/spf13/cobra"
)

///todo
//https://en.inkasex.com/videos/latest
//https://pornbox.com/application/studio/328
//https://pornbox.com/application/model/188419
///

type Host interface {
	Name() string
}

type Site struct {
	name string
	View func(string) preview.Preview
}

func (s *Site) Name() string { return s.name }

type Family struct {
	name string
	List []Site
}

func (f *Family) Name() string { return f.name }

// type Domain struct {
// 	domain string
// 	view   func(string) preview.Preview
// }

func toDomainList(name string, list []string, f func(string) preview.Preview) (d Host) {
	l := []Site{}
	for _, s := range list {
		l = append(l, Site{s, f})
	}
	d = &Family{name, l}
	return
}

// func run() (pr []preview.ContentResource, err error) {
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
// 	return
// }

var hostList = buildDomains()

func buildDomains() (d []Host) {
	d = append(d, toDomainList(`brazzer`, brazz.Domain, brazz.New))
	d = append(d, toDomainList(`pornbox`, pbox.Domain, pbox.New))
	d = append(d, toDomainList(`teamskeet`, teamskt.Domain, teamskt.New))

	d = append(d, &Site{`pornworld`, pworld.New},
		&Site{`sexmex`, smex.New},
		&Site{`vip4k`, vip.New},
		&Site{`pornhd8k`, hd8k.New},
		&Site{`nubiles`, nporn.New},
		&Site{`private`, private.New},
		&Site{`dorcel`, dorcel.New},
		&Site{`littlecaprice`, littlecaprice.New},
		&Site{`enjoyx`, enjoyx.New},
		&Site{`lifese`, lifese.New},
	)
	d = append(d, toDomainList(`familytherapyxxx`, ftherapy.Domain, ftherapy.New))
	d = append(d, toDomainList(`czechav`, czechav.Domain, czechav.New))
	d = append(d, toDomainList(`newsensations`, newsens.Domain, newsens.New))
	d = append(d, toDomainList(`allanal`, allanal.Domain, allanal.New))
	d = append(d, toDomainList(`evilangel`, evil.Domain, evil.New))
	return
}

var example = func() string {
	var buff strings.Builder
	for i, d := range hostList {
		fmt.Fprintf(&buff, "%2d: %s\n", i, d.Name())
	}
	return buff.String()
}()

func getM3U(cmd *cobra.Command, args []string, index int) (m3u string, err error) {
	if len(args) == 0 {
		err = fmt.Errorf("need at least one argument")
		return
	}

	input, err := strconv.Atoi(args[0])
	if err != nil {
		return
	}

	host := hostList[input]
	var site Site

	switch h := host.(type) {
	case *Family:
		{
			var buff strings.Builder
			for i, d := range h.List {
				fmt.Fprintf(&buff, "%2d: %s\n", i, d.Name())
			}
			cmd.Example = buff.String()

			if len(args) < 2 {
				err = fmt.Errorf("need one more argument")
				return
			}
			subIndex, err := strconv.Atoi(args[1])
			if err != nil {
				return ``, err
			}
			if subIndex < 0 || subIndex >= len(h.List) {
				err = fmt.Errorf("subIndex >= 0 && subIndex < len(f.List); subindex: %d", subIndex)
				return ``, err
			}
			site = h.List[subIndex]
		}
	case *Site:
		site = *h

	default:
		err = fmt.Errorf("switch case exhausted")
		return

	}

	pr, err := site.View(site.name).Get(index)
	if err != nil {
		return ``, err
	}
	var buffer strings.Builder
	buffer.WriteString("#EXTM3U\n")
	for _, cr := range pr {
		buffer.WriteString("#EXTINF:-1," + cr.Source + "\n" + cr.View + "\n")
	}
	m3u = buffer.String()
	return
}

func cobraPlay(cfg config.Config, str string) (err error) {
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
	_, err = stdin.Write([]byte(str))
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
			m3u, err := getM3U(cmd, args, index)
			if err != nil {
				return err
			}
			fmt.Println(m3u)
			return err
		},
	}

	play := cobra.Command{
		Use:     "play",
		Example: example,
		RunE: func(cmd *cobra.Command, args []string) error {
			m3u, err := getM3U(cmd, args, index)
			if err != nil {
				return err
			}

			return cobraPlay(cfg, m3u)
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
