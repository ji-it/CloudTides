export default function () {
    return [
        {
            name: "Dashboard",
            to: "/home",
            to1: "/manage-resources",
            htmlBefore: '<i class="icon icon-dashboard tides-fw tides-white"></i>',
            htmlAfter: "",
            show: true,
        },
        {
            name: "Contribution",
            htmlBefore: '<i class="icon icon-contribution tides-fw tides-white"></i>',
            to: "/blog-posts",
            show: false,
        },
        {
            name: "Template",
            htmlBefore: '<i class="icon icon-template tides-fw tides-white"></i>',
            to: "/add-new-post",
            show: false,
        },
        {
            name: "Settings",
            htmlBefore: '<i class="icon icon-settings tides-fw tides-white"></i>',
            to: "/components-overview",
            show: false,
        },
    ];
}
