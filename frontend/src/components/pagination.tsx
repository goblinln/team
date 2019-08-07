import * as React from 'react';

import {Icon} from './icon';
import './pagination.css';

interface PaginationProps {
    current: number;
    total: number;

    onChange?: (page: number) => void,
};

export const Pagination = (props: PaginationProps) => {
    const [current, setCurrent] = React.useState<number>(props.current);
    const [leftCtrls, setLeftCtrls] = React.useState<React.ReactNode[]>([]);
    const [midCtrls, setMidCtrls] = React.useState<React.ReactNode[]>([]);
    const [rightCtrls, setRightCtrls] = React.useState<React.ReactNode[]>([]);

    React.useEffect(() => {
        if (current < 4 || props.total == 5) {
            let left: React.ReactNode[] = [];
            for (let i = 0; i < 5 && i < props.total; i++) left.push(makeItem(i));
            setLeftCtrls(left);
        } else {
            setLeftCtrls([makeItem(0)]);
        }

        if (current+4 < props.total && current >= 4) {
            let mid: React.ReactNode[] = [];
            for (let i = 5; i > 0; i--) mid.push(makeItem(current-i+3));
            setMidCtrls(mid);
        } else {
            setMidCtrls([]);
        }

        if (current+5 > props.total && props.total > 7) {
            let right: React.ReactNode[] = [];
            for (let i = 5; i > 0; i--) right.push(makeItem(props.total-i));
            setRightCtrls(right);
        } else if (props.total > 5) {
            setRightCtrls([makeItem(props.total-1)]);
        }
    }, [props]);

    const makeItem = (idx: number) => {
        return (
            <button
                key={idx}
                className={idx==current ? 'pagination-active' : undefined} 
                onClick={() => moveTo(idx)}>
                {idx+1}
            </button>
        );
    };

    const moveTo = (page: number) => {
        setCurrent(page);
        if (props.onChange) props.onChange(page);
    };

    return (
        <div className='pagination'>
            <button key='prev' disabled={current<1} onClick={() => moveTo(current-1)}><Icon type='left'/></button>
            {leftCtrls}
            {midCtrls.length > 0 && [<button key='left-divider'>••</button>, ...midCtrls]}
            {rightCtrls.length > 0 && [<button key='right-divider'>••</button>, ...rightCtrls]}
            <button key='next' disabled={current+1>=props.total} onClick={() => moveTo(current+1)}><Icon type='right'/></button>
        </div>
    );
};