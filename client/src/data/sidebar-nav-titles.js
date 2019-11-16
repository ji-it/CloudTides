export default function () {
    return [
        {
            name: "Dashboard",
            to: "/home",
            children: [
                {
                    name: "Overview",
                    to: "/home",
                },
                {
                    name: "Manage Resources",
                    to: "/manage-resources",
                },
            ],
        },
        {
            name: "Contribution",
            to:
                "/contribution"
        }
        ,
        {
            name: "Template",
            to:
                "/add-new-post",
        }
        ,
        {
            name: "Settings",
            to:
                "/components-overview",
        }
        ,
        {
            name: "Manage Resources",
            to:
                "/manage-resources",
        }
    ]
        ;
}
