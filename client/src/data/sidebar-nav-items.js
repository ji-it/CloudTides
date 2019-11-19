export default function () {
    return [
        {
            name: "Dashboard",
            to: "/home",
            to1: "/manage-resources",
            htmlBefore: '<i class="icon tides-icon-dashboard tides-fw tides-white"></i>',
            htmlAfter: "",
            show: false,
        },
        {
            name: "Manage Resources",
            htmlBefore: '<i class="icon tides-icon-resource tides-fw tides-white"></i>',
            to: "/manage-resources",
            show: true,
        },
        {
            name: "Contribution",
            htmlBefore: '<i class="icon tides-icon-contribution tides-fw tides-white"></i>',
            to: "/contribution",
            show: false,
        },
        {
            name: "Templates",
            htmlBefore: '<i class="icon tides-icon-template tides-fw tides-white"></i>',
            to: "/templates",
            show: false,
        },
        {
            name: "Settings",
            htmlBefore: '<i class="icon tides-icon-settings tides-fw tides-white"></i>',
            to: "/settings",
            show: false,
        },
    ];
}
