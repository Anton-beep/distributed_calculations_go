import '../App.css'
import {useEffect, useState} from "react";

export const ViewExpressions = () => {
    const [expressions, setExpressions] = useState([])

    useEffect(() => {
        fetch(process.env.REACT_APP_STORAGE_API_URL + "/expression")
            .then(response => response.json())
            .then(data => setExpressions(data.expressions))
            .catch(err => console.log(err));
    }, []);


    return (
        <>
            <div className="scrollable-div">
                <h1>View All Expressions</h1>
                <ul className="list-group list-group-horizontal">
                    <li className="list-group-item">ID</li>
                    <li className="list-group-item">Value</li>
                    <li className="list-group-item">Answer</li>
                    <li className="list-group-item">Logs</li>
                    <li className="list-group-item">Ready</li>
                    <li className="list-group-item">Creation Time</li>
                    <li className="list-group-item">End Calculation Time</li>
                    <li className="list-group-item">Name of the Server</li>
                </ul>
                {expressions !== [] ? expressions.map((expression, index) =>
                    (
                        <ul className="list-group list-group-horizontal" key={index}>
                            <li className="list-group-item list-group-item-primary">{expression.id}</li>
                            <li className="list-group-item list-group-item-primary">{expression.value}</li>
                            <li className="list-group-item list-group-item-primary">{expression.answer}</li>
                            <li className="list-group-item list-group-item-primary">{expression.logs}</li>
                            <li className="list-group-item list-group-item-primary">{expression.ready}</li>
                            <li className="list-group-item list-group-item-primary">{expression.creation_time}</li>
                            <li className="list-group-item list-group-item-primary">{expression.end_calculation_time}</li>
                            <li className="list-group-item list-group-item-primary">{expression.server_name}</li>
                        </ul>
                    )
                ) : null}
            </div>
        </>
    )
}