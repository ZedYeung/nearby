import React from 'react';
import { Form, Input, Button, message } from 'antd';
import $ from 'jquery';
import { API_ROOT } from '../.env.js'
import { Link } from 'react-router-dom'

const FormItem = Form.Item;
class RegistrationForm extends React.Component {
    state = {
        confirmDirty: false
    }
    handleSubmit = (e) => {
        e.preventDefault();
        this.props.form.validateFieldsAndScroll((err, values) => {
            if (!err) {
                // console.log('Received values of form: ', values);
                $.ajax({
                    url: `${API_ROOT}/signup`,
                    method: 'POST',
                    data: JSON.stringify({
                        username: values.username,
                        password: values.password
                    })
                }).then((response) => {
                    // message.success(response);
                    this.props.history.push('/login');
                }).catch(function(error) {
                    message.error(error.responseText);
                });
            }
        });
    }
    handleConfirmBlur = (e) => {
        const value = e.target.value;
        this.setState({ confirmDirty: this.state.confirmDirty || !!value });
    }
    checkPassword = (rule, value, callback) => {
        const form = this.props.form;
        if (value && value !== form.getFieldValue('password')) {
            callback('Two passwords that you enter is inconsistent!');
        } else {
            callback();
        }
    }
    checkConfirm = (rule, value, callback) => {
        const form = this.props.form;
        if (value && this.state.confirmDirty) {
            form.validateFields(['confirm'], { force: true });
        }
        callback();
    }
    render() {
        const { getFieldDecorator } = this.props.form;
        const formItemLayout = {
            labelCol: {
                xs: { span: 24 },
                sm: { span: 6 },
            },
            wrapperCol: {
                xs: { span: 24 },
                sm: { span: 14 },
            },
        };
        const tailFormItemLayout = {
            wrapperCol: {
                xs: {
                    span: 24,
                    offset: 0,
                },
                sm: {
                    span: 14,
                    offset: 6,
                },
            },
        };
        return (
            <Form onSubmit={this.handleSubmit} className="register-form">
                <FormItem
                    {...formItemLayout}
                    label="Username"
                    hasFeedback
                >
                    {getFieldDecorator('username', {
                        rules: [{ required: true, message: 'Please input your username!', whitespace: true}],
                    })(
                        <Input />
                    )}
                </FormItem>
                <FormItem
                    {...formItemLayout}
                    label="Password"
                    hasFeedback
                >
                    {getFieldDecorator('password', {
                        rules: [{
                            required: true, message: '6-16 characters, at least one uppercase letter, one lowercase letter, one number',
                            pattern: /^(?=.*[a-z])(?=.*[A-Z])(?=.*\d)[a-zA-Z\d]{6,16}$/
                        }, {
                            validator: this.checkConfirm,
                        }],
                    })(
                        <Input type="password" />
                    )}
                </FormItem>
                <FormItem
                    {...formItemLayout}
                    label="Confirm Password"
                    hasFeedback
                >
                    {getFieldDecorator('confirm', {
                        rules: [{
                            required: true, message: 'Please confirm your password!',
                        }, {
                            validator: this.checkPassword,
                        }],
                    })(
                        <Input type="password" onBlur={this.handleConfirmBlur} />
                    )}
                </FormItem>
                <FormItem {...tailFormItemLayout}>
                    <Button type="primary" htmlType="submit">Register</Button>
                    <p>I already have an account, go back to <Link to='/Login'>login</Link></p>
                </FormItem>
            </Form>
        );
    }
}

export const Register = Form.create()(RegistrationForm);