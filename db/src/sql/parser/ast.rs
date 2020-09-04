use crate::sql::parser::lexer::Token;

pub enum Statement {
    Select {
        select: Vec<Token>,
        from: Token,
        r#where: Option<Token>,
    },
}

impl Statement {
    pub fn build_select(query: String) -> Statement::Select {

    }
}