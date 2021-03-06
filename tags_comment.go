package pongo2

type tagCommentNode struct{}

func (node *tagCommentNode) Execute(ctx *ExecutionContext) (string, error) {
	return "", nil
}

func tagCommentParser(doc *Parser, start *Token, arguments *Parser) (INodeTag, error) {
	comment_node := &tagCommentNode{}

	_, err := doc.WrapUntilTag("endcomment")
	if err != nil {
		return nil, err
	}

	if arguments.Count() != 0 {
		return nil, arguments.Error("Tag 'comment' does not take any argument.", nil)
	}

	return comment_node, nil
}

func init() {
	RegisterTag("comment", tagCommentParser)
}
