import { Link } from 'react-router-dom';
import '../App.css'

export const Header = () => {
    return (
        <header className="header">
            <Link to="/" className="btn btn-primary">Home</Link>
            <Link to="/inputExpression" className="btn btn-primary">Input New Expression</Link>
            <Link to="/viewExpressions" className="btn btn-primary">View All Expressions</Link>
            <Link to="/operations" className="btn btn-primary">View Operations And Execution Times</Link>
            <Link to="/computingPowers" className="btn btn-primary">View Computing Powers</Link>
        </header>
    )
}