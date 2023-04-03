import React, { useEffect, useRef, useState } from 'react';
import { PlusOutlined } from '@ant-design/icons';
import type { InputRef } from 'antd';
import { Input, Tag } from 'antd';
import { TweenOneGroup } from 'rc-tween-one';


interface State {
    inputVisible: boolean;
    inputValue: string;
    tags: string[];
}

interface Props {
    record: {
        tags: string[]
    },
    onChange: (tags: string[]) => void;
}

export const HomeTag: React.FC<Props> = ({ record, onChange }) => {

    const [inputVisible, setInputVisible] = useState(false);
    const [inputValue, setInputValue] = useState('');
    const [tags, setTags] = useState(record.tags || []);

    const inputRef = useRef<InputRef>(null);


    useEffect(() => {
        if (record.tags !== undefined) {
            setTags(record.tags)
        }
    }, []);

    const handleClose = (removedTag: string) => {
        const newTags = tags.filter(tag => tag !== removedTag);
        console.log(newTags);
        setTags(newTags)
        onChange(newTags)
    };

    const showInput = () => {
        if (!inputVisible) {
            setInputVisible(true);
            setTimeout(() => {
                inputRef.current?.focus();
            }, 200);
        }
    };

    const handleInputChange = (e: React.ChangeEvent<HTMLInputElement>) => {
        setInputValue(e.target.value)
    };

    const handleInputConfirm = () => {
        if (inputValue && tags.indexOf(inputValue) === -1) {
            const newtags = [...tags]; 
            newtags.push(inputValue)
            console.info("new tags", newtags)
            setTags(newtags)
            onChange(newtags)
        }

        setInputValue("")
        setInputVisible(false)
    };

    const forMap = (tag: string) => {
        const tagElem = (
            <Tag
                closable
                onClose={e => {
                    e.preventDefault();
                    handleClose(tag);
                }}
            >
                {tag}
            </Tag>
        );
        return (
            <span key={tag} style={{ display: 'inline-block' }}>
                {tagElem}
            </span>
        );
    };

    const tagChild = tags.map(forMap);

    return (
        <>
            <div style={{ marginBottom: 16 }}>
                <TweenOneGroup
                    enter={{
                        scale: 0.8,
                        opacity: 0,
                        type: 'from',
                        duration: 100,
                    }}
                    onEnd={e => {
                        if (e.type === 'appear' || e.type === 'enter') {
                            (e.target as any).style = 'display: inline-block';
                        }
                    }}
                    leave={{ opacity: 0, width: 0, scale: 0, duration: 200 }}
                    appear={false}
                >
                    {tagChild}
                </TweenOneGroup>
            </div>
            {
                inputVisible && (
                    <Input
                        type="text"
                        size="small"
                        style={{ width: 78 }}
                        value={inputValue}
                        onChange={handleInputChange}
                        onBlur={handleInputConfirm}
                        onPressEnter={handleInputConfirm}
                        ref={inputRef}
                    />
                )
            }
            {
                !inputVisible && (
                    <Tag onClick={showInput} className="site-tag-plus">
                        <PlusOutlined /> New Tag
                    </Tag>
                )
            }
        </>
    );

};
