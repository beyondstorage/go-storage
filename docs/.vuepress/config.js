module.exports = {
    title: 'storage',
    description: 'An application-oriented unified storage layer for Golang.',
    smoothScroll: true,
    plugins: [
        '@vuepress/last-updated',
        '@vuepress/active-header-links',
        '@vuepress/google-analytics',
        {
            'ga': 'UA-51515330-5'
        }
    ],
    locales: {
        '/': {
            lang: 'en-US',
            title: 'storage',
            description: 'An application-oriented unified storage layer for Golang.'
        },
        '/zh-CN/': {
            lang: 'zh-CN',
            title: 'storage',
            description: '面向应用的 Golang 抽象存储层'
        }
    },
    themeConfig: {
        nav: [
            {text: 'Services', link: '/services/'},
            {text: 'Operations', link: '/operations/'},
            {text: 'Design', link: '/design/'},
            {text: 'Spec', link: '/spec/'}
        ],
        sidebar: "auto",
        repo: 'Xuanwo/storage',
        docsDir: 'docs',
        editLinks: true,
        locales: {
            "/": {
                selectText: 'Languages',
                label: 'English',
                ariaLabel: 'Languages',
                editLinkText: 'Edit this page on GitHub',
                serviceWorker: {
                    updatePopup: {
                        message: "New content is available.",
                        buttonText: "Refresh"
                    }
                },
                nav: [
                    {text: 'Services', link: '/services/', ariaLabel: 'Services'},
                    {text: 'Operations', link: '/operations/', ariaLabel: 'Operations'},
                    {text: 'Design', link: '/design/', ariaLabel: 'Design'},
                    {text: 'Spec', link: '/spec/', ariaLabel: 'Spec'}
                ]
            },
            "/zh-CN/": {
                selectText: '语言',
                label: '简体中文',
                ariaLabel: '语言',
                editLinkText: '在 Github 上编辑此页',
                serviceWorker: {
                    updatePopup: {
                        message: "发现新内容可用",
                        buttonText: "刷新"
                    }
                },
                nav: [
                    {text: '服务', link: '/zh-CN/services/'},
                    {text: '操作', link: '/zh-CN//operations/'},
                    {text: '设计', link: '/zh-CN/design/'},
                    {text: '规范', link: '/zh-CN/spec/'}
                ]
            },
        }
    },
}
