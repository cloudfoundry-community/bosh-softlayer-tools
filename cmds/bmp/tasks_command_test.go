package bmp_test

import (
	"errors"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	fakes "github.com/cloudfoundry-community/bosh-softlayer-tools/clients/fakes"
	cmds "github.com/cloudfoundry-community/bosh-softlayer-tools/cmds"
	bmp "github.com/cloudfoundry-community/bosh-softlayer-tools/cmds/bmp"
)

var _ = Describe("BMP: tasks command", func() {
	var (
		args    []string
		options cmds.Options
		cmd     cmds.Command

		fakeBmpClient *fakes.FakeBmpClient
	)

	BeforeEach(func() {
		args = []string{"bmp", "tasks"}
		options = cmds.Options{
			Verbose: false,
			Latest:  0,
		}

		fakeBmpClient = fakes.NewFakeBmpClient("fake-username", "fake-password", "http://fake.url.com", "fake-config-path")
		cmd = bmp.NewTasksCommand(options, fakeBmpClient)
	})

	Describe("NewTasksCommand", func() {
		It("create new TasksCommand", func() {
			Expect(cmd).ToNot(BeNil())

			cmd2 := bmp.NewTasksCommand(options, fakeBmpClient)
			Expect(cmd2).ToNot(BeNil())
			Expect(cmd2).To(Equal(cmd))
		})
	})

	Describe("#Name", func() {
		It("returns the name of a TasksCommand", func() {
			Expect(cmd.Name()).To(Equal("tasks"))
		})
	})

	Describe("#Description", func() {
		It("returns the description of a TasksCommand", func() {
			Expect(cmd.Description()).To(Equal(`List tasks: \"option --latest number\", \"The number of latest tasks, default is 50\"`))
		})
	})

	Describe("#Usage", func() {
		It("returns the usage text of a TasksCommand", func() {
			Expect(cmd.Usage()).To(Equal("bmp tasks"))
		})
	})

	Describe("#Options", func() {
		It("returns the options of a TasksCommand", func() {
			Expect(cmds.EqualOptions(cmd.Options(), options)).To(BeTrue())
		})
	})

	Describe("#Validate", func() {
		It("validates a good TasksCommand", func() {
			validate, err := cmd.Validate()
			Expect(validate).To(BeTrue())
			Expect(err).ToNot(HaveOccurred())
		})
	})

	Describe("#Execute", func() {
		Context("executes a good TasksCommand", func() {
			BeforeEach(func() {
				fakeBmpClient.TasksResponse.Status = 200
				fakeBmpClient.TasksErr = nil
				options = cmds.Options{
					Verbose: false,
					Latest:  1,
				}
			})

			It("executes a good TasksCommand without specifying latest", func() {
				rc, err := cmd.Execute([]string{"bmp", "tasks"})
				Expect(rc).To(Equal(0))
				Expect(err).ToNot(HaveOccurred())
			})

			It("executes a good TasksCommand with specifying latest", func() {
				cmd = bmp.NewTasksCommand(options, fakeBmpClient)

				rc, err := cmd.Execute([]string{"bmp", "tasks", "--latest=1"})
				Expect(rc).To(Equal(0))
				Expect(err).ToNot(HaveOccurred())
			})
		})

		Context("executes a bad TasksCommand", func() {
			Context("executes TasksCommand with error", func() {
				BeforeEach(func() {
					fakeBmpClient.TasksResponse.Status = 500
					fakeBmpClient.TasksErr = errors.New("500")
				})

				It("executes with error", func() {
					rc, err := cmd.Execute([]string{"bmp", "tasks"})
					Expect(rc).To(Equal(1))
					Expect(err).To(HaveOccurred())
				})
			})

			Context("TasksCommand response different than 200", func() {
				BeforeEach(func() {
					fakeBmpClient.TasksResponse.Status = 404
				})

				It("response code different than 200", func() {
					rc, err := cmd.Execute([]string{"bmp", "tasks"})
					Expect(rc).To(Equal(404))
					Expect(err).ToNot(HaveOccurred())
				})
			})
		})

	})
})
