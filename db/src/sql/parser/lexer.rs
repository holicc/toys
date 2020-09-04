use std::iter::Peekable;
use std::option::Option::Some;
use std::str::Chars;

use crate::error::{Error, Result};
use crate::sql::parser::keyword::Keyword;

#[derive(Debug, Clone, PartialEq)]
pub enum Token {
    Number(String),
    String(String),
    Keyword(Keyword),
    Ident(String),
}

impl From<Keyword> for Token {
    fn from(keyword: Keyword) -> Self {
        Self::Keyword(keyword)
    }
}

pub struct Lexer<'a> {
    iter: Peekable<Chars<'a>>
}

impl<'a> Iterator for Lexer<'a> {
    type Item = Result<Token>;

    fn next(&mut self) -> Option<Self::Item> {
        match self.scan() {
            Ok(Some(token)) => Some(Ok(token)),
            Ok(None) => match self.iter.peek() {
                Some(c) => Some(Err(Error::Parse(format!("Unexpected character {}", c)))),
                None => None,
            }
            Err(err) => Some(Err(err)),
        }
    }
}

impl<'a> Lexer<'a> {
    pub fn new(input: &'a str) -> Lexer<'a> {
        Lexer { iter: input.chars().peekable() }
    }

    fn scan(&mut self) -> Result<Option<Token>> {
        //skip whitespace
        self.skipp_whitespace();
        //scan to token
        match self.iter.peek() {
            Some('\'') => self.scan_string(),
            Some(c) if c.is_digit(10) => Ok(self.scan_num()),
            Some(c) if c.is_alphabetic() => Ok(self.scan_keyword()),
            Some(_) => self.scan_symbol(),
            None => Ok(None),
        }
    }

    fn skipp_whitespace(&mut self) {
        while let Some(c) = self.iter.peek() {
            if c.is_whitespace() {
                self.iter.next();
            } else {
                break;
            }
        }
    }

    fn scan_string(&mut self) -> Result<Option<Token>> {
        self.iter.next();
        let mut r = String::new();
        while let Some(c) = self.iter.next() {
            if c == '\'' { break; }
            r.push(c);
        }
        Ok(Some(Token::String(r)))
    }

    fn scan_num(&mut self) -> Option<Token> {
        let a: String = self.iter.peek()
            .iter()
            .map(|c| **c)
            .take_while(|x|
                (*x) == '+' || (*x) == '-' || (*x) == '.' || x.is_digit(10)
            )
            .collect();

        Some(Token::Number(a))
    }

    fn scan_symbol(&mut self) -> Result<Option<Token>> {
        unimplemented!()
    }

    fn scan_keyword(&mut self) -> Option<Token> {
        let mut a = String::new();
        //
        while let Some(c) = self.iter.peek() {
            if !c.is_alphabetic() {
                break;
            }
            a.push(self.iter.next().unwrap());
        }
        //
        Keyword::from_str(a.as_str())
            .map(Token::Keyword)
            .or_else(|| Some(Token::Ident(a.to_lowercase())))
    }
}

