package nav

import "fmt"

type Link struct {
	Title string
	URL   string

	idx      int
	active   bool
	parent   *Link
	children []Link
}

func (l *Link) Add(link Link) *Link {
	nl := Link{
		Title:  link.Title,
		URL:    link.URL,
		parent: l,
		idx:    len(l.children),
	}

	for _, cl := range link.Children() {
		nl.Add(cl)
	}

	l.children = append(l.children, nl)

	return l
}

func (l *Link) Activate(url string) bool {
	for childIdx := range l.children {
		if l.children[childIdx].Activate(url) {
			l.active = true
		}
	}

	if l.URL == url {
		l.active = true
	}

	return l.active
}

func (l Link) Active() bool {
	return l.active
}

func (l Link) Children() []Link {
	return l.children
}

func (l Link) Depth() int {
	if l.HasParent() {
		return l.parent.Depth() + 1
	}

	return 0
}

func (l Link) HasParent() bool {
	return l.parent != nil
}

func (l Link) Index() int {
	return l.idx
}

func (l *Link) Ref() string {
	ref := "ref"

	if l.HasParent() {
		ref = l.parent.Ref()
	}

	return fmt.Sprintf("%s-%d", ref, l.idx)
}
