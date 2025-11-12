The query package provides a parser for predicates based on the [ANTLR](https://www.antlr.org/) grammar defined in
Query.g4. The parser returns a tree representation of the predicate which is then visited to produce the final
`Predicate` structure.

The ANTLR grammar is generated using the antlr4-go tool. To regenerate the grammar, install Antlr and run
the target `generate-query-grammar` in `common/Makefile`.

NOTE If the grammar changes, the visitor logic in `predicateBuilder.go` will need to be updated.