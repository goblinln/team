import * as React from 'react';
import * as ReactDOM from 'react-dom';
import * as ReactMarkdown from 'react-markdown';

import {
    Button,
    Icon,
    Input,
    Popover,
    Row,
    Upload,
    message,
} from 'antd';

/**
 * 默认样式
 */
import './Markdown.css';

/**
 * Markdown渲染组件可配置参数
 */
export interface IRendererProps {
    /**
     * Markdown内容
     */
    source?: string;
}

/**
 * 渲染组件
 */
export const Renderer = ((props: IRendererProps) => {
    /**
     * 打开图片预览器
     */
    const openImgViewer = (src: string) => {
        let layer = document.createElement('div');
        layer.id = 'markdown-image-viewer';
        layer.style.position = 'absolute';
        layer.style.left = '0';
        layer.style.right = '0';
        layer.style.top = '0';
        layer.style.bottom = '0';
        layer.style.zIndex = '10000';
        layer.style.backgroundColor = 'rgba(0, 0, 0, .5)';
        layer.style.display = 'flex';
        layer.style.flexDirection = 'row';
        layer.style.justifyContent = 'center';
        layer.style.alignItems = 'center';
        document.body.appendChild(layer);

        let close = () => {
            ReactDOM.render(null, layer);
            layer.remove();
        };

        layer.addEventListener('click', () => close());

        ReactDOM.render((
            <div style={{position: 'relative', backgroundColor: 'white', padding: 8}}>
                <Button shape='circle' icon='close' style={{position: 'absolute', right: -16, top: -16, width: 32, height: 32}} onClick={() => close()}/>
                <div style={{maxWidth: 'calc(100vw - 80px)', maxHeight: 'calc(100vh - 80px)', overflow: 'auto'}}>
                    <img src={encodeURI(src)}/>
                </div>
            </div>            
        ), layer);
    }

    /**
     * 对React-mardown的渲染逻辑扩展一下
     */
    const customRenderers = {
        listItem: (props: any) => {
            let liProps = props['data-sourcepos'] ? { 'data-sourcepos': props['data-sourcepos'] } : {};

            if (props.checked === null) {
                return <li {...liProps}>{props.children}</li>;
            } else {
                return (
                    <li className='task-list-item' {...liProps}>
                        <input type='checkbox' checked={props.checked} readOnly={true}/>
                        {props.children}
                    </li>
                );
            }
        },
        image: (props: any) => {
            return <img {...props} onClick={ev => openImgViewer(ev.currentTarget.src)}/>
        }
    };

    return (
        <ReactMarkdown 
            source={props.source || ''} 
            escapeHtml={false} 
            className='markdown-body' 
            renderers={customRenderers}/>
    )
});

/**
 * Markdown编辑器的可配置属性
 */
export interface IEditorProps {
    /**
     * 行数
     */
    rows?: number;
    /**
     * 不使用row，直接设置高度
     */
    height?: number | string;
    /**
     * 编辑内容
     */
    value?: string;
    /**
     * 内容发生变化时回调
     */
    onChange?: (data: string) => void;
    /**
     * 自定义上传功能，done必需调用！
     */
    onUpload?: (image: File, done: (url: string) => void) => void;
}

/**
 * Markdown编辑器
 */
export const Editor = (props: IEditorProps) => {
    
    /**
     * 状态列表
     */
    const [content, setContent] = React.useState<string>(props.value || '');
    const [showUploader, setShowUploader] = React.useState<boolean>(false);
    const [showPreview, setShowPreview] = React.useState<boolean>(false);
    const textArea = React.useRef<any>();

    /**
     * 图片工具
     */
    const imageButton = (
        <Popover
            title='图片工具'
            content={<ImageHelper onUpload={props.onUpload} onDone={url => insertImage(url)}/>}
            visible={showUploader}
            trigger='click'
            onVisibleChange={() => setShowUploader(prev => !prev)}>
            <Button title='上传图片' style={{padding: '0 10px'}}><Icon type='file-image'/></Button>
        </Popover>
    );

    /**
     * 菜单栏配置
     */
    const toolbar = [
        [
            {useText: true, caption: 'H1', tooltip: '一级标题', modify: (data: string) => `  \n# ${data || '一级标题'}  \n`},
            {useText: true, caption: 'H2', tooltip: '二级标题', modify: (data: string) => `  \n## ${data || '二级标题'}  \n`},
            {useText: true, caption: 'H3', tooltip: '三级标题', modify: (data: string) => `  \n### ${data || '三级标题'}  \n`},
        ],
        [
            {caption: 'bold', tooltip: '粗体', modify: (data: string) => `**${data || '粗体'}**`},
            {caption: 'italic', tooltip: '斜体', modify: (data: string) => `*${data || '斜体'}*`},
            {caption: 'strikethrough', tooltip: '删除线', modify: (data: string) => `<s>${data || '删除线'}</s>`},
            {caption: 'underline', tooltip: '下划线', modify: (data: string) => `<u>${data || '下划线'}</u>`},
        ],
        [
            {caption: 'ordered-list', tooltip: '有序表', modify: (data: string) => `  \n1. ${data || '第一项'}  \n2. 第二项  \n`},
            {caption: 'unordered-list', tooltip: '无序表', modify: (data: string) => `  \n* 第一项 ${data || '第一项'}  \n* 第二项  \n`},
            {caption: 'schedule', tooltip: '任务列表', modify: (data: string) => `  \n* [ ] ${data || '第一项'}  \n* [ ] 第二项  \n`},
        ],
        [
            {caption: 'code', tooltip: '代码', modify: (data: string) => `  \n\`\`\`\n${data || '代码'}\n\`\`\`  \n`},
            {caption: 'form', tooltip: '引用', modify: (data: string) => `  \n> ${data || '引用内容'}\n\n`},
            {caption: 'table', tooltip: '表格', modify: (data: string) => `\n\n| 标题1 | 标题2 |\n|---|---|  \n${data || ''}`},
            {caption: 'link', tooltip: '链接', modify: (data: string) => `[显示内容](${data || '连接地址'})`},
            {element: imageButton},
        ]
    ];

    /**
     * 编辑文件变化时
     */
    const onContentChange = (ev: React.ChangeEvent<HTMLTextAreaElement>) => {
        setContent(ev.target.value);
        props.onChange && props.onChange(ev.target.value);
    };

    /**
     * 根据光标位置和选中文本，编辑文本
     */
    const modifyContent = (action: (data: string) => string) => {
        let elem = (textArea.current.textAreaRef as HTMLTextAreaElement);
        let data = [content];

        let start = elem.selectionStart;
        let end = elem.selectionEnd;
        data = [
            content.substr(0, start),
            content.substr(start, end - start),
            content.substr(end),
        ];

        let modified = (data[0] || '') + action(data[1]) + (data[2] || '');
        setContent(modified);
        props.onChange && props.onChange(modified);

        elem.focus();
    };

    /**
     * 插入图片功能
     */
    const insertImage = (url: string) => {
        if (url && url.length > 0) modifyContent(data => `![${data || '输入图片TIPS'}](${url})`);
        setShowUploader(false);
    };

    /**
     * 粘贴图片功能
     */
    const pasteImage = (ev: React.ClipboardEvent<HTMLTextAreaElement>) => {
        if (!props.onUpload) return;

        if (!(ev.clipboardData && ev.clipboardData.items && ev.clipboardData.items.length > 0)) {
            return;
        }

        let item = ev.clipboardData.items[0];
        if (item.kind == 'file' && item.type.indexOf('image') != -1) {
            let file = item.getAsFile();
            props.onUpload(file, url => insertImage(url));
        }
    };
    
    return (
        <div style={{padding: '0px 4px'}}>
            <Row type='flex' style={{marginBottom: 4}}>
                {toolbar.map((group: any[], i: number) => {
                    return (
                        <Button.Group key={i} style={{marginRight: 4}}>
                            {group.map((opt, j) => {
                                if (opt.element) return opt.element;

                                return (
                                    <Button key={`${i}_${j}`} title={opt.tooltip} style={{padding: '0 10px'}} onClick={() => modifyContent(opt.modify)}>
                                        {opt.useText ? opt.caption : <Icon type={opt.caption} />}
                                    </Button>
                                );
                            })}
                        </Button.Group>
                    );
                })}

                <Button.Group>
                    <Button 
                        title='预览' 
                        type={showPreview ? 'primary' : 'default'}
                        style={{padding: '0 10px'}}
                        onClick={() => setShowPreview(prev => !prev)}>
                        <Icon type='eye' />
                    </Button>
                </Button.Group>
            </Row>

            <Row>
                {showPreview ? (
                    <div style={{height: (props.rows ? props.rows * 22 : (props.height ? props.height : '100%')), padding: 4, overflowY: 'auto', border: '1px solid rgba(0, 0, 0, .15)'}}>
                        <Renderer source={content}/>
                    </div>
                ) : (
                    <Input.TextArea
                        ref={textArea}
                        rows={props.rows || null}
                        style={{height: props.height || null}}
                        value={content}
                        onChange={onContentChange}
                        onPaste={pasteImage} />
                )}
            </Row>
        </div>
    );
};

/**
 * 上传工具的可配置属性
 */
interface IImageHelperProps {
    /**
     * 启用“插入图片功能”需要配置的自定义上传功能。
     */
    onUpload?: (image: File, done: (url: string) => void) => void;
    /**
     * 完成上传或输入后的回调
     */
    onDone?: (url: string) => void;
}

/**
 * 用于编辑器的图片上传工具
 */
const ImageHelper = (props: IImageHelperProps) => {

    /**
     * 状态列表
     */
    const [url, setUrl] = React.useState<string>(null);

    /**
     * 启用自定义上传
     */
    const customUpload = (image: File) => {
        if (props.onUpload == null) {
            message.error('当前编辑器未配置上传功能！', 1);
        } else {
            props.onUpload(image, path => setUrl(path));
        }

        return false;
    };

    return (
        <div>
            <Row>
                <Input
                    addonBefore={<Upload name='image' showUploadList={false} accept='image/*' beforeUpload={customUpload}><Button type='link' size='small'>上传图片</Button></Upload>}
                    placeholder='或直接输入图片URL'
                    value={url}
                    onChange={(data) => setUrl(data.target.value)}/>
            </Row>

            <Row style={{marginTop: 8}}>
                <Button type='primary' onClick={() => props.onDone && props.onDone(url)} block>插入图片</Button>
            </Row>
        </div>
    );
};