import '../App.css'
import {Link} from "react-router-dom";

export const Home = () => {
    return (
        <>
            <h1>Home</h1>
            <p>This project assumes all standard mathematical operations (+, /, *, -) need a lot of time to be
                calculated. Therefore it would be logical to create a system which will organize a work of several
                machines to calculated given expressions as fast as possible.</p>
            <p>There are 5 main functionalities to control the process:</p>
            <ul>
                <li>
                    <Link to="/inputExpression">Input an expression</Link>
                </li>
                <li>
                    <Link to="/viewExpressions">View all expressions</Link>
                </li>
                <li>
                    <Link to="/operations">View operations and execution times</Link>
                </li>
                <li>
                    <Link to="/computingPowers">View computing powers</Link>
                </li>
            </ul>
            <p>See more detailed description on the <a href="https://github.com/Anton-beep/distributed_calculations_go">GitHub repository</a>
            </p>
        </>
    )
}