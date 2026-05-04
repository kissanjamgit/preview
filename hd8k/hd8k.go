package hd8k

import (
	"fmt"
	"os"
	"regexp"

	"github.com/kissanjamgit/ext"
	"resty.dev/v3"
)

// https://en8.pornhd8k.me/movies/new-sofie-marie-double-penetration-spit-roasting-oh-my-bully-and-nerd-share-holes-23-04-2026-anal-doublepenetration-hardcore-milf-bigtits-threesome-family-stepmom-iluvy-luluvdo-com-doodstream-com
// https://cdn-aws-exp.cdnamz.me/playlists/e34f082c-fd2e-4ba2-a648-f9610d539ec2/u8lmahqjne.m3u8?md5=0OZ1p5_nbJfyYsZmt8FYUQ&expires=1776956326

type HD8k struct {
	Source string
}

func (h *HD8k) Resource(client *resty.Client) (cr ext.ContentResource, err error) {
	f, err := os.ReadFile("hd8k/content.txt")
	if err != nil {
		return
	}
	submatch := regexp.MustCompile(`<input id="uuid" type="hidden" value="([^"]+)"`).FindStringSubmatch(string(f))
	if len(submatch) < 2 {
		err = fmt.Errorf("len(submatch) < 2")
		return
	}

	fmt.Println(submatch[1])
	panic("unimplemented")
}

func (h *HD8k) Download(cr ext.ContentResource) (err error) {
	panic("unimplemented")
}

func New(source string) ext.Site {
	return &HD8k{Source: source}
}
