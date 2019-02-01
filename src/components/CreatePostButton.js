import React from 'react';
import { Modal, Button, message } from 'antd';
import {WrappedCreatePostForm} from './CreatePostForm';
import {API_ROOT, POS_KEY, AUTH_PREFIX, TOKEN_KEY} from '../.env';
import $ from 'jquery';

export class CreatePostButton extends React.Component {
    state = {
        visible: false,
        confirmLoading: false,
    }

    showModal = () => {
        this.setState({
            visible: true,
        });
    }

    handleOk = () => {
        this.form.validateFields((error, values) =>{
            if(!error){
                console.log('Received values of form: ', values);

                const {lat, lon} = JSON.parse(localStorage.getItem(POS_KEY));
                // const {lat, lon} = {lat:47.6492236871, lon:-122.3937977};
                const formData = new FormData();
                // jitter
                formData.set('lat', lat + Math.random() * 0.1 - 0.05);
                formData.set('lon', lon + Math.random() * 0.1 - 0.05);
                formData.set('message', values.message);

                values.images.forEach((image) => {
                    formData.append('images[]', image.originFileObj, image.name);
                })

                this.setState({confirmLoading: true,});
                $.ajax({
                    url:`${API_ROOT}/post`,
                    method: 'POST',
                    headers: {
                        Authorization: `${AUTH_PREFIX} ${localStorage.getItem(TOKEN_KEY)}`,
                    },
                    processData: false,
                    contentType: false,
                    dataType: 'text',
                    data: formData,
                }).then(() => {
                    message.success('created a post successfully.');
                    this.form.resetFields();//ant design method
                }, (error)=>{
                    message.error(error.responseText);
                    this.form.resetFields();//ant design method
                }).then(()=>{
                    this.props.loadNearByPosts().then(()=>{
                        this.setState({ visible: false, confirmLoading: false,});
                    });
                }).catch((error)=>{
                    message.error('create post failed.');
                    console.log(error);
                });
            }
        });
    }
    handleCancel = () => {
        console.log('Clicked cancel button');
        this.setState({
            visible: false,
        });
    }

    saveFormRef = (form) => {
        this.form = form;
    }

    render() {
        const { visible, confirmLoading } = this.state;
        return (
            <div>
                <Button type="primary" onClick={this.showModal}>Create New Post</Button>
                <Modal title="Create New Post"
                       visible={visible}
                       onOk={this.handleOk}
                       okText="Create"
                       confirmLoading={confirmLoading}
                       onCancel={this.handleCancel}
                       cancelText="Cancel"
                >
                    <WrappedCreatePostForm ref={this.saveFormRef}/>
                </Modal>
            </div>
        );
    }
}