module.exports = {
    title: 'storage',
    description: 'An application-oriented unified storage layer for Golang.',
    smoothScroll: true,
    sidebar: [
        '/',
        '/page-a',
        ['/page-b', 'Explicit link text']
    ],
    plugins: [
        '@vuepress/last-updated',
        '@vuepress/active-header-links',
        '@vuepress/google-analytics',
        {
            'ga': 'UA-51515330-5'
        }
    ],
    themeConfig: {
        nav: [
            {text: 'Services', link: '/services/'},
            {text: 'Design', link: '/design/'},
            {text: 'Spec', link: '/spec/'}
        ],
        sidebar: "auto",
        repo: 'Xuanwo/storage',
        docsDir: 'docs',
        editLinks: true,
    },
}