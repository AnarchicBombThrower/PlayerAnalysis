import "./header.css"

function Header(){
    return(
        <header className="siteHeader">
            <h1>Player Analysis</h1>
            <nav>
                <ul>
                    <li><a href="https://github.com/AnarchicBombThrower">Check out my github!</a></li>
                </ul>
            </nav>
        </header>
    );
}

export default Header