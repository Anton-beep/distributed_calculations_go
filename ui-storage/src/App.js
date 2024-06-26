import {Header} from './components/Header';
import {Home} from './components/Home';
import {InputExpression} from './components/InputExpression';
import {ViewExpressions} from "./components/ViewExpressions";
import {Operations} from "./components/Operations";
import {ComputingPowers} from "./components/ComputingPowers";
import {Footer} from "./components/Footer";
import {Register} from "./components/Registration";
import {Login} from "./components/Login";
import {Logout} from "./components/Logout";
import {Profile} from "./components/Profile";

import {Route, Routes} from 'react-router-dom';


function App() {
    return (
        <>
            <link href="https://cdn.jsdelivr.net/npm/bootstrap@5.2.3/dist/css/bootstrap.min.css" rel="stylesheet"
                  integrity="sha384-rbsA2VBKQhggwzxH7pPCaAqO46MgnOM80zW1RWuH61DGLwZJEdK2Kadq2F9CUG65"
                  crossOrigin="anonymous"/>
            <script src="https://cdn.jsdelivr.net/npm/bootstrap@5.2.3/dist/js/bootstrap.bundle.min.js"
                    integrity="sha384-kenU1KFdBIe4zVF0s0G1M5b4hcpxyD9F7jL+jjXkk+Q2h455rYXK/7HAuoJl+0I4"
                    crossOrigin="anonymous"></script>
            <Header/>
            <main>
                <Routes>
                    <Route path="/" element={<Home/>}/>
                    <Route path="/register" element={<Register/>}/>
                    <Route path="/login" element={<Login/>}/>
                    <Route path="/profile" element={<Profile/>}/>
                    <Route path="/logout" element={<Logout/>}/>
                    <Route path="/inputExpression" element={<InputExpression/>}/>
                    <Route path="/viewExpressions" element={<ViewExpressions/>}/>
                    <Route path="/operations" element={<Operations/>}/>
                    <Route path="/computingPowers" element={<ComputingPowers/>}/>
                </Routes>
            </main>
            <Footer/>
        </>
    )
}

export default App;
