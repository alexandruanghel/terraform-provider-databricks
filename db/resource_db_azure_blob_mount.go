package db

import (
	"fmt"
	"github.com/databrickslabs/databricks-terraform/client/service"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/helper/validation"
	"strings"
)

func resourceAzureBlobMount() *schema.Resource {
	return &schema.Resource{
		Create: resourceAzureBlobMountCreate,
		Read:   resourceAzureBlobMountRead,
		Delete: resourceAzureBlobMountDelete,

		Schema: map[string]*schema.Schema{
			"cluster_id": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"container_name": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"storage_account_name": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"directory": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
				//Default:  "/",
				ForceNew: true,
				ValidateFunc: func(val interface{}, key string) (warns []string, errors []error) {
					directory := val.(string)
					if strings.HasPrefix(directory, "/") {
						return
					} else {
						errors = append(errors, fmt.Errorf("%s must start with /, got: %s", key, val))
					}
					return
				},
			},
			"mount_name": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"auth_type": {
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validation.StringInSlice([]string{"SAS", "ACCESS_KEY"}, false),
			},
			"token_secret_scope": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"token_secret_key": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
		},
	}
}

func resourceAzureBlobMountCreate(d *schema.ResourceData, m interface{}) error {
	client := m.(service.DBApiClient)
	clusterId := d.Get("cluster_id").(string)
	err := changeClusterIntoRunningState(clusterId, client)
	if err != nil {
		return err
	}
	containerName := d.Get("container_name").(string)
	storageAccountName := d.Get("storage_account_name").(string)
	directory := d.Get("directory").(string)
	mountName := d.Get("mount_name").(string)
	authType := d.Get("auth_type").(string)
	tokenSecretScope := d.Get("token_secret_scope").(string)
	tokenSecretKey := d.Get("token_secret_key").(string)

	blobMount := service.NewAzureBlobMount(containerName, storageAccountName, directory, mountName, authType,
		tokenSecretScope, tokenSecretKey)

	err = blobMount.Create(client, clusterId)

	d.SetId(mountName)

	err = d.Set("cluster_id", clusterId)
	if err != nil {
		return err
	}
	err = d.Set("mount_name", mountName)
	if err != nil {
		return err
	}
	err = d.Set("auth_type", authType)
	if err != nil {
		return err
	}
	err = d.Set("token_secret_scope", tokenSecretScope)
	if err != nil {
		return err
	}
	err = d.Set("token_secret_key", tokenSecretKey)
	if err != nil {
		return err
	}

	return resourceAzureBlobMountRead(d, m)
}
func resourceAzureBlobMountRead(d *schema.ResourceData, m interface{}) error {
	client := m.(service.DBApiClient)
	clusterId := d.Get("cluster_id").(string)
	err := changeClusterIntoRunningState(clusterId, client)
	if err != nil {
		return err
	}
	containerName := d.Get("container_name").(string)
	storageAccountName := d.Get("storage_account_name").(string)
	directory := d.Get("directory").(string)
	mountName := d.Get("mount_name").(string)
	authType := d.Get("auth_type").(string)
	tokenSecretScope := d.Get("token_secret_scope").(string)
	tokenSecretKey := d.Get("token_secret_key").(string)

	blobMount := service.NewAzureBlobMount(containerName, storageAccountName, directory, mountName, authType,
		tokenSecretScope, tokenSecretKey)

	url, err := blobMount.Read(client, clusterId)
	if err != nil {
		//Reset id in case of inability to find mount
		if strings.Contains(err.Error(), "Unable to find mount point!") {
			d.SetId("")
			return nil
		}
		return err
	}
	container, storageAcc, dir, err := service.ProcessAzureWasbAbfssUris(url)
	err = d.Set("container_name", container)
	if err != nil {
		return err
	}
	err = d.Set("storage_account_name", storageAcc)
	if err != nil {
		return err
	}
	err = d.Set("directory", dir)
	return err
}

func resourceAzureBlobMountDelete(d *schema.ResourceData, m interface{}) error {
	client := m.(service.DBApiClient)
	clusterId := d.Get("cluster_id").(string)
	err := changeClusterIntoRunningState(clusterId, client)
	if err != nil {
		return err
	}
	containerName := d.Get("container_name").(string)
	storageAccountName := d.Get("storage_account_name").(string)
	directory := d.Get("directory").(string)
	mountName := d.Get("mount_name").(string)
	authType := d.Get("auth_type").(string)
	tokenSecretScope := d.Get("token_secret_scope").(string)
	tokenSecretKey := d.Get("token_secret_key").(string)

	blobMount := service.NewAzureBlobMount(containerName, storageAccountName, directory, mountName, authType,
		tokenSecretScope, tokenSecretKey)
	return blobMount.Delete(client, clusterId)
}