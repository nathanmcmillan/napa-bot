import React from 'react';
import ReactDOM from 'react-dom';
import Hello from './hello.js';
import World from './world.js';

class Root extends React.Component {
    render() {
        return <h1>i am groot the root</h1>
    }
}

ReactDOM.render(<Root />, document.getElementById('root'));
ReactDOM.render(<Hello />, document.getElementById('hello'));
ReactDOM.render(<World />, document.getElementById('world'));
console.log('is this on?');