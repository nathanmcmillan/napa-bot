import React, { Component } from 'react';
import { View, Text } from 'react-native';

export default class App extends Component {
    render() {
        return (
            <View style={styles.container}>
                <Text style={styles.intro}>Hello world!</Text>
            </View>
        );
    }
}