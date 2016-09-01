//Package seq_definition is where the regexes for file sequences are defined, you may edit this file if needed
//Original regexes are:
//Reducer: ".*([\.\_\ \/\\]([0-9]+)\.)\w{2,4}$"
//Expander: ".*[\.\_\ \/\\](\[[0-9-,]+\])\."
//Note that there is a group inside the regex around the number component of the file or sequence listing.
package seq_definition

type Seq_definition struct {
	ReducerRegex  string
	ExpanderRegex string
}

func SeqDefinition() (Seq_definition, error) {
	seq_def := Seq_definition{
		ReducerRegex:  `.*(([\.\_\ \/\\])([0-9]+)\.(\w{2,4}$))`,
		ExpanderRegex: `.*[\.\_\ \/\\](\[[0-9-,]+\])\.`,
	}
	return seq_def, nil
}
