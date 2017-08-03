function ContainerDetailsViewModel(data) {
  this.Id = data.Id;
  this.State = data.State;
  this.Created = data.Created;
  this.Name = data.Name;
  this.NetworkSettings = data.NetworkSettings;
  this.Args = data.Args;
  this.Image = data.Image;
  this.Config = data.Config;
  this.HostConfig = data.HostConfig;
  if (data.dockm) {
    if (data.dockm.ResourceControl) {
      this.ResourceControl = new ResourceControlViewModel(data.dockm.ResourceControl);
    }
  }
}
