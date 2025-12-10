/*
 * Copyright (c) 2025 Metaform Systems, Inc
 *
 * This program and the accompanying materials are made available under the
 * terms of the Apache License, Version 2.0 which is available at
 * https://www.apache.org/licenses/LICENSE-2.0
 *
 * SPDX-License-Identifier: Apache-2.0
 *
 * Contributors:
 *      Metaform Systems, Inc. - initial API and implementation
 */

grammar Query;

// Defines the grammar for CFM queries

predicate
    : compoundPredicate EOF
    | atomicPredicate EOF
    | TRUE_KEYWORD EOF
    ;

compoundPredicate
    : atomicPredicate ((AND | OR) atomicPredicate)+
    | '(' compoundPredicate ')'
    ;

atomicPredicate
    : fieldPath comparisonOperator value
    | fieldPath IS NULL_KEYWORD
    | fieldPath IS NOT NULL_KEYWORD
    | fieldPath IN LPAREN valueList RPAREN
    | fieldPath NOT IN LPAREN valueList RPAREN
    ;

fieldPath
    : IDENTIFIER (DOT IDENTIFIER)*
    ;

comparisonOperator
    : EQ
    | NE
    | GT
    | GTE
    | LT
    | LTE
    | LIKE
    | NOT_LIKE
    | CONTAINS
    | STARTS_WITH
    | ENDS_WITH
    ;

value
    : stringValue
    | numericValue
    | booleanValue
    | NULL_KEYWORD
    ;

valueList
    : value (COMMA value)*
    ;

stringValue
    : STRING
    ;

numericValue
    : NUMBER
    ;

booleanValue
    : TRUE_KEYWORD
    | FALSE_KEYWORD
    ;

// Lexer Rules

// Keywords
AND: ('AND' | 'and');
OR: ('OR' | 'or');
IS: ('IS' | 'is');
NOT: ('NOT' | 'not');
IN: ('IN' | 'in');
NULL_KEYWORD: ('NULL' | 'null');
TRUE_KEYWORD: ('TRUE' | 'true');
FALSE_KEYWORD: ('FALSE' | 'false');
LIKE: ('LIKE' | 'like');
NOT_LIKE: ('NOT' | 'not') WS? ('LIKE' | 'like');
CONTAINS: ('CONTAINS' | 'contains');
STARTS_WITH: ('STARTS_WITH' | 'starts_with');
ENDS_WITH: ('ENDS_WITH' | 'ends_with');

// Operators
EQ: '=';
NE: '!=';
GT: '>';
GTE: '>=';
LT: '<';
LTE: '<=';

// Punctuation
LPAREN: '(';
RPAREN: ')';
DOT: '.';
COMMA: ',';

// Literals
STRING
    : '\'' (~['\\\r\n] | '\\' .)* '\''
    ;

NUMBER
    : '-'? DIGIT+ ('.' DIGIT+)?
    ;

DIGIT
    : [0-9]
    ;

IDENTIFIER
    : [a-zA-Z_] [a-zA-Z0-9_]*
    ;

// Skip whitespace
WS
    : [ \t\n\r]+ -> skip
    ;