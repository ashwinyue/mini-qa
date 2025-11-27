/**
 * 建议问题列表组件
 * 
 * 渲染可点击的建议问题按钮
 */
import React from 'react'
import { Button } from 'antd'

/**
 * 建议问题列表组件
 * 
 * @param {Object} props
 * @param {Array} props.suggestions - 建议问题数组
 * @param {Function} props.onClick - 点击回调
 */
const SuggestionList = ({ suggestions = [], onClick }) => {
    if (!suggestions || suggestions.length === 0) {
        return null
    }

    const handleClick = (suggestion) => {
        onClick && onClick(suggestion)
    }

    return (
        <div className="flex flex-wrap gap-2 mt-2">
            {suggestions.map((suggestion, index) => (
                <Button
                    key={index}
                    size="small"
                    type="default"
                    onClick={() => handleClick(suggestion)}
                    className="text-blue-500 border-blue-300 hover:border-blue-500"
                >
                    {suggestion}
                </Button>
            ))}
        </div>
    )
}

export default SuggestionList
