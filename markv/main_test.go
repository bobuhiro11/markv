package main

import (
	"testing"
)

func TestConvert(t *testing.T) {
	input := `# h1 Heading
## h2 Heading
### h3 Heading
#### h4 Heading
##### h5 Heading
###### h6 Heading

Horizontal Rules
______________

Sample Image
![Sample Image](http://www.ess.ic.kanagawa-it.ac.jp/std_img/colorimage/Lenna.jpg)

 H1  | H2   | H3
-----|------|-----
row1 | item | item
row2 | item | item
row3 | item | item

- Link: [some link](http://example.com/sample.html)
- **Strong**
- ` + "code" + `
- itemA
- itemB
- itemC
    + itemD

<span>span</span>

<div>
div
</div>

1. enum1
  - A
  - B
2. enum2

> block quote


` + "```ruby" + `
a = "Hello, World"
puts foo
` + "```" + "\n```go" + `
func sum(a, b int) bool {
  return a + b
}
` + "```"
	expected := `h1Heading
 h2Heading
  h3Heading
   h4Heading
    h5Heading
     h6Heading
HorizontalRules
♦━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━♦
SampleImage
⡇⡇⢅⢕⢑⢕⡑⣕⡑⡅⡣⡣⡫⣪⠪⢈
⡂⡇⠢⢢⢑⠔⡕⡵⡽⣦⠪⢌⣪⠨⢐⢔
⡂⡇⠅⢕⢔⠱⡱⡙⠮⡺⡽⡵⡹⢐⢜⢜
⡂⡇⠅⡅⠇⢕⢐⢬⢪⢳⡣⠅⡂⡕⣕⢕
⠂⡇⢕⢌⠪⡂⡒⢕⢜⢜⢔⢑⢌⢎⢎⢮
⠡⡣⡑⡸⡰⡊⠌⠜⡌⡎⡂⠢⡣⡣⡽⡽
⢪⢪⠂⡊⢎⢆⢅⠱⡘⡜⣦⢱⢱⢨⡎⢕
⢸⢸⠐⢌⢜⠐⢄⠕⡅⡣⡫⣇⠅⢇⢑⠔

+------+------+------+
|  H1  |  H2  |  H3  |
+------+------+------+
| row1 | item | item |
| row2 | item | item |
| row3 | item | item |
+------+------+------+
  - Link:somelink
  - Strong
  - code
  - itemA
  - itemB
  - itemC
    - itemD
<span>span</span>
<div>
div
</div>
  1. enum1
    - A
    - B
  2. enum2

⦀ blockquote

+--------------------+
| a = "Hello, World" |
| puts foo           |
+--------------------+
+---------------------------+
| func sum(a, b int) bool { |
|   return a + b }          |
+---------------------------+

`
	actual := convert(input)
	if actual != expected {
		t.Fatalf("expected = %x\nactual = %x\n", expected, actual)
	}
}
