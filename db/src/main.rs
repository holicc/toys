use std::io;
use std::io::Write;
use std::process::exit;

use db::sql::executor::Executor;

enum FLAG {
    EXIT,
    HELP,
    QUERY(String),
}

struct Repl {
    i: io::Stdin,
    o: io::Stdout,
}

fn main() {
    //
    let exec = Executor::new();
    //
    let mut repl = Repl::new();
    loop {
        repl.print("db:> ");
        match repl.read_line() {
            FLAG::EXIT => {
                repl.print("Bye!");
                exit(0);
            }
            FLAG::HELP => {
                repl.print("This is fake database not need help! \n");
            }
            FLAG::QUERY(query) => {
                exec.execute(query)
            }
        }
    }
}


impl Repl {
    fn new() -> Repl {
        return Repl {
            i: io::stdin(),
            o: io::stdout(),
        };
    }

    fn print(&mut self, s: &str) {
        self.o.write(s.as_bytes()).unwrap();
        self.o.flush().unwrap();
    }

    fn read_line(&self) -> FLAG {
        let mut r = String::new();
        self.i.read_line(&mut r).unwrap();
        match r.trim_end_matches(|x| x == '\n') {
            "\\q" => FLAG::EXIT,
            "\\h" => FLAG::HELP,
            query => FLAG::QUERY(String::from(query).to_lowercase())
        }
    }
}

