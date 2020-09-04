use crate::error::{Result,Error};
use crate::sql::parser::ast::Statement;
use crate::sql::parser::lexer::Lexer;

pub struct Executor {}

struct ResultSet {}

impl Executor {
    pub fn new() -> Executor {
        Executor {}
    }

    pub fn execute(self, query: String) -> Result<ResultSet> {
        let mut lex = Lexer::new(query.into());
        let statement = match lex.scan() {
            Ok(Some(token)) => Statement::build(token),
            Ok(None) => return Err(Error::Parse("Nothing to do".into())),
            Err(e) => return Err(e),
        };
        //
        engine
    }
}