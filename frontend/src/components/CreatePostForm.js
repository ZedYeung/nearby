import React, { Component } from 'react';
import {Form, Input, Upload, Icon} from 'antd';

const FormItem = Form.Item;

class CreatePostForm extends Component {
    normFile = (e) => {
        console.log('Upload event:', e);
        if (Array.isArray(e)) {
            console.log("e is array")
            return e;
        }
        return e && e.fileList;
    }

    beforeUpload = () =>{// Uploading will be stopped with false or a rejected Promise
        return false;
    }

    render(){
        // this.props.form is WrappedCreatePostForm
        const { getFieldDecorator } = this.props.form;
        const formItemLayout = {
            labelCol: {span: 6},
            wrapperCol: {span: 14},
        };
        return (
            <Form layout="vertical">
                <FormItem
                    {...formItemLayout}
                    label="Message">
                    {getFieldDecorator('message', {
                        rules: [{ required: true, message: 'Please input a message.' }],
                    })(
                        <Input />
                    )}
                </FormItem>
                <FormItem
                    {...formItemLayout}
                    label="Image"
                >
                    <div className="dropbox">
                        {getFieldDecorator('images', {
                            valuePropName: 'fileList',
                            getValueFromEvent: this.normFile,
                            rules: [{ required: true, message: 'Please select an image.' }],
                        })(
                            <Upload.Dragger
                                name="images"
                                multiple={true}
                                beforeUpload={this.beforeUpload}
                                accept="image/*"
                            >
                                <p className="ant-upload-drag-icon">
                                    <Icon type="inbox" />
                                </p>
                                <p className="ant-upload-text">Click or drag file to this area to upload</p>
                                <p className="ant-upload-hint">Support for a single or bulk upload.</p>
                            </Upload.Dragger>
                        )}
                    </div>
                </FormItem>
            </Form>
        );
    }
}

export const WrappedCreatePostForm = Form.create()(CreatePostForm);