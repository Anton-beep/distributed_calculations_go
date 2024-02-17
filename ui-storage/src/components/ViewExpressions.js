import '../App.css'
import {useEffect, useState} from "react";

export const ViewExpressions = () => {
    const [expressions, setExpressions] = useState([])

    useEffect(() => {
        let addr;
        if (process.env.REACT_APP_STORAGE_API_URL === undefined) {
            addr = process.env.REACT_APP_STORAGE_API_URL + "/expression"
        } else {
            addr = "http://localhost:8080/api/v1/expression"
        }
        fetch(addr)
            .then(response => response.json())
            .then(data => {
                data.expressions.sort((a, b) => (a.id > b.id) ? -1 : 1)
                setExpressions(data.expressions)
            })
            .catch(err => {
                console.log(err)
            });
    }, []);

    const showReady = (ready) => {
        switch (ready) {
            case 0:
                return (<>
                    <div>
                        {"Expression is waiting to be calculated"}
                    </div>
                    <img src={process.env.PUBLIC_URL + '/clock-history.svg'} alt="clock" width="32" height="32"/>
                </>)
            case 1:
                return (<>
                    <div>
                        {"Server is calculating this expression"}
                    </div>
                    <img src={process.env.PUBLIC_URL + '/calculator.svg'} alt="calculator" width="32" height="32"/>
                </>)
            case 2:
                return (<>
                    <div>
                        {"Expression is calculated"}
                    </div>
                    <img src={process.env.PUBLIC_URL + '/check2-circle.svg'} alt="check" width="32" height="32"/>
                </>)
            case 3:
                return (<>
                    <div>
                        {"Error calculating expression, see logs"}
                    </div>
                    <img src={process.env.PUBLIC_URL + '/exclamation-octagon.svg'} alt="error" width="32" height="32"/>
                </>)
            default:
                return (<>
                    <div>
                        {"Unknown status"}
                    </div>
                    <img src={process.env.PUBLIC_URL + '/exclamation-octagon.svg'} alt="error" width="32" height="32"/>
                </>)
        }
    }


    return (
        <>
            <div className="scrollable-div">
                <h1>View All Expressions</h1>
                <table className="table table-striped-columns">
                    <thead>
                    <tr>
                        <th>ID</th>
                        <th>Value</th>
                        <th>Answer</th>
                        <th>Logs (output from calculation server)</th>
                        <th>Status</th>
                        <th>Creation Time</th>
                        <th>End Calculation Time</th>
                        <th>Name of the Server</th>
                    </tr>
                    </thead>
                    <tbody>
                    {expressions !== [] ? expressions.map((expression, index) =>
                        (
                            <tr key={index}>
                                <th>{expression.id}</th>
                                <th>{expression.value}</th>
                                <th>{expression.answer}</th>
                                <th>{expression.logs.split("\n").map((el, key) => {
                                    return <div key={key}>{el}</div>;
                                })}</th>
                                <th>{showReady(expression.ready)}</th>
                                <th>{expression.creation_time}</th>
                                <th>{expression.end_calculation_time}</th>
                                <th>{expression.server_name}</th>
                            </tr>
                        )
                    ) : null}
                    </tbody>
                </table>
            </div>
        </>
    )
}