use std::{collections::HashMap, io::{self, Read}};

const SYNTAX_CHECK_ENABLED: bool = false;

macro_rules! expect_token {
    ($inter:expr, $x:expr) => {
        if SYNTAX_CHECK_ENABLED && Some($x) != $inter.current(){
            unreachable!("Expected: {:?} but got: {:?}", $x, $inter.current());
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

struct Interpreter{
    tokens: Vec<Token>,
    val_indices: Vec<usize>,
    val_index: usize,
    vars: Vec<usize>,
    index: usize,
}

impl Interpreter {
    pub fn new(tokens: Vec<Token>, val_indices: Vec<usize>, vars: Vec<usize>) -> Self{
        Self{
            tokens,
            val_indices,
            val_index: 0,
            vars,
            index: 0,
        }
    }

    pub fn interpret(&mut self) {
        self.interpret_statements();
        for (i, v) in self.vars.iter().enumerate(){
            println!("{}: {}", i, v);
        }
    }

    fn interpret_statements(&mut self){
        while let Some(tok) = self.current() && tok != Token::End && tok != Token::Eof{
            self.interpret_statement();
        }
    }

    fn interpret_statement(&mut self){
        let current = self.current();
        match current {
            Some(Token::Loop) | Some(Token::While) => {
                self.index += 1;

                expect_token!(self, Token::Ident);
                self.index += 1;

                let loop_val_index = self.val_index;
                self.val_index += 1;

                expect_token!(self, Token::Do);
                self.index += 1;

                let loop_amount = self.vars[self.val_indices[loop_val_index]];
                if loop_amount == 0{
                    self.jump_to_end();
                    return;
                }
                let current_index = self.index;
                let current_val_index = self.val_index;
                if let Some(Token::Loop) = current{
                    for _ in 0..loop_amount{
                        self.index = current_index;
                        self.val_index = current_val_index;
                        self.interpret_statements();
                    }
                }
                if let Some(Token::While) = current{
                    while self.vars[self.val_indices[loop_val_index]] > 0{
                        self.index = current_index;
                        self.val_index = current_val_index;
                        self.interpret_statements();
                    }
                }

                expect_token!(self, Token::End);
                self.index += 1;
            }
            Some(Token::Ident) => {
                self.index += 1;

                let current_ident_index = self.val_index;
                self.val_index += 1;

                expect_token!(self, Token::Assign);
                self.index += 1;

                expect_token!(self, Token::Ident);
                self.index += 1;
                let other_ident_index = self.val_index;
                self.val_index += 1;

                let op_tok = self.current();
                let is_add = if let Some(Token::Plus) = op_tok {true} else if let Some(Token::Minus) = op_tok {false} else {unreachable!("{op_tok:?}")};
                self.index += 1;

                expect_token!(self, Token::Num);
                self.index += 1;
                let num_index = self.val_index;
                self.val_index += 1;
                let number = self.vars[self.val_indices[num_index]];

                let other_val = self.vars[self.val_indices[other_ident_index]];

                /*println!("{}={}, {}={}, {}, {}, at {}",
                    self.val_indices[current_ident_index],
                    self.vars[self.val_indices[current_ident_index]],
                    self.val_indices[other_ident_index],
                    self.vars[self.val_indices[other_ident_index]],
                    is_add,
                    number,
                    self.index
                );*/

                if is_add{
                    self.vars[self.val_indices[current_ident_index]] = other_val.checked_add(number).expect("Overflowed add")
                } else {
                    self.vars[self.val_indices[current_ident_index]] = other_val.checked_sub(number).unwrap_or(0)
                }
            }
            tok => unreachable!("Expected LOOP, WHILE or IDENT but got: {:?} at {}", tok, self.index)
        }
    }

    fn jump_to_end(&mut self){
        let mut count = 1;
        let mut j = self.index;
        while count > 0{
            let tok = self.tokens.get(j);
            match tok {
                Some(Token::Loop) | Some(Token::While) => count += 1,
                Some(Token::End) => count -= 1,
                Some(Token::Ident) | Some(Token::Num) => self.val_index += 1,
                _ => {}
            }
            j += 1;
        }
        self.index = j;
    }

    fn current(&self) -> Option<Token>{
        self.tokens.get(self.index).copied()
    }
}

fn main() -> io::Result<()>{
    let mut input = String::new();
    let _ = io::stdin().read_to_string(&mut input)?;
    let (tokens, var_indices, vars) = lex(input.as_str());
    println!("{:?}", tokens);
    let mut interpreter = Interpreter::new(tokens, var_indices, vars);
    interpreter.interpret();
    Ok(())
}
