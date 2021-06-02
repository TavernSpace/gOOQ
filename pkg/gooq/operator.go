package gooq

type Operator string

var (
	OperatorNil = Operator("")

	OperatorIsNull    = Operator("IS NULL")
	OperatorIsNotNull = Operator("IS NOT NULL")

	OperatorAsc  = Operator("ASC")
	OperatorDesc = Operator("DESC")

	// logical operators
	// https://www.postgresql.org/docs/11/functions-logical.html
	OperatorAnd = Operator("AND")
	OperatorOr  = Operator("OR")
	OperatorNot = Operator("NOT")

	// comparison functions and operators
	// https://www.postgresql.org/docs/11/functions-comparison.html
	OperatorLt    = Operator("<")
	OperatorLte   = Operator("<=")
	OperatorGt    = Operator(">")
	OperatorGte   = Operator(">=")
	OperatorEq    = Operator("=")
	OperatorNotEq = Operator("!=")
	OperatorLike  = Operator("LIKE")
	OperatorILike = Operator("ILIKE")
	OperatorIsDistinctFrom = Operator("IS DISTINCT FROM")

	// mathematical operators
	// https://www.postgresql.org/docs/11/functions-math.html
	// [Good First Issue][Help Wanted] TODO: implement remaining
	OperatorAdd  = Operator("+")
	OperatorSub  = Operator("-")
	OperatorMult = Operator("*")
	OperatorDiv  = Operator("/")
	OperatorSqrt = Operator("|/")

	// Table 9.13. Bit String Operators
	// https://www.postgresql.org/docs/11/functions-bitstring.html
	// [Good First Issue][Help Wanted] TODO: implement remaining

	// 9.7.1. LIKE
	// 9.7.2. SIMILAR TO Regular Expressions
	// 9.7.3. POSIX Regular Expressions
	// https://www.postgresql.org/docs/11/functions-matching.html
	// [Good First Issue][Help Wanted] TODO: implement remaining

	// Array Comparisons
	// https://www.postgresql.org/docs/11/functions-comparisons.html
	OperatorIn    = Operator("IN")
	OperatorNotIn = Operator("NOT IN")
)

func (op Operator) String() string {
	return string(op)
}
