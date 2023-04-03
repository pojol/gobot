const DocsPage = () => {

  let iframeStyle = {
    position: "absolute",
    width: "100%",
    height: "1024px",
    top: 0,
    border: "0",
    frameborder:"no"
  }

  let rootStyle = {
    position: "relative",
    width: "100%",
    border: "2px black solid;"
  }

  return (
    <div style={rootStyle as React.CSSProperties}>
      <iframe src="https://pojol.gitee.io/gobot/#/" style={iframeStyle as React.CSSProperties}>
      </iframe>
    </div>
  );
};

export default DocsPage;