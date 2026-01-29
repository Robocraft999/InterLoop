use std::collections::HashMap;
use std::io;
use std::io::Read;

macro_rules! expect_token {
    ($parser:expr, $x:expr) => {
        if Some($x) != $parser.current(){
            unreachable!("Expected: {:?} but got: {:?}", $x, $parser.current());
        }
    };
}

#[derive(Debug, PartialEq, Clone, Copy)]
enum Token{
    Eof,
    Plus,
    Minus,
    Assign,
    Loop,
    While,
    Do,
    End,
    Ident,
    Num,
}

fn lex(input: &str) -> (Vec<Token>, Vec<usize>, Vec<usize>){
    let mut vars = Vec::new();
    let mut var_indices = Vec::new();
    let mut var_mapping = HashMap::new();
    let tokens = input.split_whitespace().map(|raw|{
        match raw{
            "" => Token::Eof,
            "+" => Token::Plus,
            "-" => Token::Minus,
            ":=" => Token::Assign,
            "LOOP" => Token::Loop,
            "WHILE" => Token::While,
            "DO" => Token::Do,
            "END" => Token::End,
            num_or_ident => {
                let var_count = var_mapping.len();
                let optional_value = usize::from_str_radix(num_or_ident, 10);
                let index = var_mapping.entry(num_or_ident).or_insert_with(|| {
                    vars.push(optional_value.clone().unwrap_or(0));
                    var_count
                });
                var_indices.push(*index);
                optional_value.map(|_| Token::Num).unwrap_or(Token::Ident)
            }
        }
    }).chain(vec![Token::Eof]).collect();
    println!("{:?}", var_mapping);
    (tokens, var_indices, vars)
}

type Value = usize;
type VarIndex = usize;

#[derive(Debug, PartialEq, Clone, Copy)]
enum Statement{
    // amount, end_index
    Loop(VarIndex, usize),
    // var_index, end_index
    While(VarIndex, usize),
    // target, source, constant, type (false = add)
    SingleAssignment(VarIndex, VarIndex, Value, bool),
    // target, source, type (false = add)
    LoopAssignment(VarIndex, VarIndex, bool),
    // target, constant
    ConstAssignment(VarIndex, Value),
    // target, factor1, factor2
    Multiplication(VarIndex, VarIndex, VarIndex),
    // condition, else_index
    If(VarIndex, usize),
}
struct Parser {
    tokens: Vec<Token>,
    var_indices: Vec<VarIndex>,
    var_index: usize,
    vars: Vec<Value>,
    statements: Vec<Statement>,
    index: usize,
}

impl Parser {
    fn new(tokens: Vec<Token>, var_indices: Vec<VarIndex>, vars: &Vec<Value>) -> Self {
        Self{
            tokens,
            var_indices,
            var_index: 0,
            vars: vars.clone(),
            statements: Vec::new(),
            index: 0,
        }
    }

    fn parse(mut self) -> Vec<Statement> {
        self.parse_statements();
        self.statements
    }

    fn parse_statements(&mut self) {
        while let Some(tok) = self.current() && tok != Token::End && tok != Token::Eof{
            self.parse_statement();
        }
    }

    fn parse_statement(&mut self) {
        let tok = self.current();
        match tok{
            Some(Token::Loop) => {
                self.index += 1;

                expect_token!(self, Token::Ident);
                self.index += 1;
                let loop_var_index = self.var_indices[self.var_index];
                self.var_index += 1;

                expect_token!(self, Token::Do);
                self.index += 1;

                self.statements.push(Statement::Loop(0, 0));

                let pre_statement_count = self.statements.len();
                self.parse_statements();
                let post_statement_count = self.statements.len();

                expect_token!(self, Token::End);
                self.index += 1;

                let is_single_statement = post_statement_count - pre_statement_count == 1;

                if is_single_statement &&
                    let Some(Statement::SingleAssignment(target, source, constant, typ)) = self.statements.last().cloned() &&
                    target == source &&
                    constant == 1 //TODO make it work for any number
                {
                    self.statements.pop().unwrap(); // single assign
                    self.statements.pop().unwrap(); // loop placeholder
                    self.statements.push(Statement::LoopAssignment(target, loop_var_index, typ));
                } else if is_single_statement &&
                    let Some(Statement::LoopAssignment(target, source, typ)) = self.statements.last().cloned() &&
                    typ == false
                {
                    self.statements.pop().unwrap(); // loop assign
                    self.statements.pop().unwrap(); // loop placeholder
                    self.statements.push(Statement::Multiplication(target, loop_var_index, source));
                } else {
                    self.statements[pre_statement_count-1] = Statement::Loop(loop_var_index, post_statement_count);
                }
            }
            Some(Token::While) => {
                self.index += 1;

                expect_token!(self, Token::Ident);
                self.index += 1;
                let loop_var_index = self.var_indices[self.var_index];
                self.var_index += 1;

                expect_token!(self, Token::Do);
                self.index += 1;

                self.statements.push(Statement::While(0, 0));

                let pre_statement_count = self.statements.len();
                self.parse_statements();
                let post_statement_count = self.statements.len();

                expect_token!(self, Token::End);
                self.index += 1;


                if self.contains_var_reset(loop_var_index, &self.statements[pre_statement_count..]){
                    self.statements[pre_statement_count-1] = Statement::If(loop_var_index, post_statement_count);
                } else {
                    self.statements[pre_statement_count-1] = Statement::While(loop_var_index, post_statement_count);
                }
            }
            Some(Token::Ident) => {
                self.index += 1;
                let target_index = self.var_indices[self.var_index];
                self.var_index += 1;

                expect_token!(self, Token::Assign);
                self.index += 1;

                expect_token!(self, Token::Ident);
                self.index += 1;
                let source_index = self.var_indices[self.var_index];
                self.var_index += 1;

                let typ = match self.current(){
                    Some(Token::Plus) => false,
                    Some(Token::Minus) => true,
                    _ => unreachable!()
                };
                self.index += 1;

                expect_token!(self, Token::Num);
                self.index += 1;
                let constant = self.vars[self.var_indices[self.var_index]];
                self.var_index += 1;

                if self.is_zero_reg(source_index) && !typ && source_index != target_index{
                    self.statements.push(Statement::ConstAssignment(target_index, constant))
                } else {
                    self.statements.push(Statement::SingleAssignment(target_index, source_index, constant, typ));
                }
            }
            token => unreachable!("{:?} cannot be used as statement start", token),
        }
    }

    fn current(&self) -> Option<Token>{
        self.tokens.get(self.index).copied()
    }

    fn is_zero_reg(&self, reg_index: VarIndex) -> bool {
        self.statements
            .iter()
            .filter(|s| if let Statement::SingleAssignment(reg, _, _, _) = s && reg == &reg_index {true} else {false})
            .count() == 0
    }

    fn contains_var_reset(&self, var_index: VarIndex, stmts: &[Statement]) -> bool {
        stmts
            .iter()
            .filter(|s|
                if let Statement::ConstAssignment(target, constant) = s &&
                    target == &var_index
                {true} else {false}
            )
            .count() == 1
    }
}

fn interpret(statements: &[Statement], vars: &mut Vec<Value>, start_index: usize) {
    let mut index = start_index;
    let statement_amount = statements.len();
    while index < statement_amount {
        //println!("{}{}: {:?}", " ".repeat(start_index), index, statements[index]);
        match statements[index] {
            Statement::Loop(var_index, end_index) => {
                let amount = vars[var_index];
                if amount == 0{
                    index = end_index;
                    continue;
                }
                for _ in 0..amount {
                    interpret(&statements[0..end_index], vars, index + 1);
                }
                index = end_index;
            }
            Statement::While(var_index, end_index) => {
                let amount = vars[var_index];
                if amount == 0{
                    index = end_index;
                    continue;
                }
                while vars[var_index] > 0{
                    interpret(&statements[0..end_index], vars, index + 1);
                }
                index = end_index;
            }
            Statement::SingleAssignment(target, source, constant, typ) => {
                if !typ{
                    vars[target] = vars[source] + constant;
                } else {
                    vars[target] = vars[source].checked_sub(constant).unwrap_or(0);
                }
                index += 1;
            }
            Statement::LoopAssignment(target, source, typ) => {
                if !typ{
                    vars[target] += vars[source];
                } else {
                    vars[target] = vars[target].checked_sub(vars[source]).unwrap_or(0);
                }
                index += 1;
            }
            Statement::ConstAssignment(target, constant) => {
                vars[target] = constant;
                index += 1;
            }
            Statement::Multiplication(target, factor1, factor2) => {
                vars[target] += vars[factor1] * vars[factor2];
                index += 1;
            }
            Statement::If(cond, else_index) => {
                if vars[cond] == 0 {
                    index = else_index;
                    continue
                }
                index += 1;
            }
        }
    }
}

fn main() -> io::Result<()>{
    let mut input = String::new();
    let _ = io::stdin().read_to_string(&mut input)?;
    let (tokens, var_indices, mut vars) = lex(input.as_str());
    //println!("{:?}", tokens);
    let mut parser = Parser::new(tokens, var_indices, &vars);
    let statements = parser.parse();
    println!("{:?}", statements);
    interpret(statements.as_slice(), &mut vars, 0);
    for (i, v) in vars.iter().enumerate(){
        println!("{}: {}", i, v);
    }
    Ok(())
}
