#[cfg(test)]
mod tests {
    use crate::sql::parser::keyword::Keyword;
    use crate::sql::parser::lexer::{Lexer, Token};

    #[test]
    fn scan_string() {
        let mut lex = Lexer::new(" 'select' 'insert'  'hello world'");
        let select = lex.next().unwrap().unwrap();
        let insert = lex.next().unwrap().unwrap();
        let hello = lex.next().unwrap().unwrap();
        //
        assert_eq!(Token::String(String::from("select")), select);
        assert_eq!(Token::String(String::from("insert")), insert);
        assert_eq!(Token::String(String::from("hello world")), hello);
    }

    #[test]
    fn scan_keyword() {
        let mut lex = Lexer::new("  SELECT INSERT   ");
        let select = lex.next().unwrap().unwrap();
        let insert = lex.next().unwrap().unwrap();
        //
        assert_eq!(Token::Keyword(Keyword::from_str("select").unwrap()), select);
        assert_eq!(Token::Keyword(Keyword::from_str("insert").unwrap()), insert);
    }
}

