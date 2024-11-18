import * as blogPluginExports from '@docusaurus/plugin-content-blog';

const defaultBlogPlugin = blogPluginExports.default;

async function blogPluginExtended(...pluginArgs) {
  const blogPluginInstance = await defaultBlogPlugin(...pluginArgs);

  const pluginOptions = pluginArgs[1];

  return {
    ...blogPluginInstance,
    contentLoaded: async function (params) {
      const {content, actions} = params;

      const recentPostsLimit = 3;
      const recentPosts = [...content.blogPosts].splice(0, recentPostsLimit);

      async function createRecentPostModule(blogPost, index) {
        return {
          metadata: await actions.createData(
            `home-page-recent-post-metadata-${index}.json`,
            JSON.stringify({
              title: blogPost.metadata.title,
              description: blogPost.metadata.description,
              frontMatter: blogPost.metadata.frontMatter,
            }),
          ),

          Preview: {
            __import: true,
            path: blogPost.metadata.source,
            query: {
              truncated: true,
            },
          },
        };
      }

      actions.addRoute({
        path: '/',
        exact: true,

        component: '@site/src/components/Home/index.tsx',

        modules: {
          homePageBlogMetadata: await actions.createData(
            'home-page-blog-metadata.json',
            JSON.stringify({
              blogTitle: pluginOptions.blogTitle,
              blogDescription: pluginOptions.blogDescription,
              totalPosts: content.blogPosts.length,
              totalRecentPosts: recentPosts.length,
            }),
          ),
          recentPosts: await Promise.all(
            recentPosts.map(createRecentPostModule),
          ),
        },
      });

      return blogPluginInstance.contentLoaded(params);
    },
  };
}

module.exports = {
  ...blogPluginExports,
  default: blogPluginExtended,
};
